package service

import (
	"context"
	"fmt"
	"time"

	service_workspace_statuses "getsturdy.com/api/pkg/workspaces/statuses/service"

	"getsturdy.com/api/pkg/activity"
	"getsturdy.com/api/pkg/activity/sender"
	service_activity "getsturdy.com/api/pkg/activity/service"
	"getsturdy.com/api/pkg/analytics"
	service_analytics "getsturdy.com/api/pkg/analytics/service"
	"getsturdy.com/api/pkg/changes"
	"getsturdy.com/api/pkg/changes/message"
	service_changes "getsturdy.com/api/pkg/changes/service"
	workers_ci "getsturdy.com/api/pkg/ci/workers"
	service_codebase "getsturdy.com/api/pkg/codebases/service"
	service_comments "getsturdy.com/api/pkg/comments/service"
	"getsturdy.com/api/pkg/events"
	eventsv2 "getsturdy.com/api/pkg/events/v2"
	"getsturdy.com/api/pkg/snapshots"
	service_snapshots "getsturdy.com/api/pkg/snapshots/service"
	worker_snapshots "getsturdy.com/api/pkg/snapshots/worker"
	service_users "getsturdy.com/api/pkg/users/service"
	service_view "getsturdy.com/api/pkg/view/service"
	vcs_view "getsturdy.com/api/pkg/view/vcs"
	"getsturdy.com/api/pkg/workspaces"
	db_workspaces "getsturdy.com/api/pkg/workspaces/db"
	service_workspaces "getsturdy.com/api/pkg/workspaces/service"
	"getsturdy.com/api/vcs"
	"getsturdy.com/api/vcs/executor"

	git "github.com/libgit2/git2go/v33"
	"go.uber.org/zap"
)

var (
	ErrNotAllowedUnhealthyWorkspace = fmt.Errorf("not allowed to land workspace, it has unhealthy statuses")
)

type Service struct {
	logger *zap.Logger

	workspaceWriter db_workspaces.WorkspaceWriter

	usersService             service_users.Service
	workspaceService         *service_workspaces.Service
	changeService            *service_changes.Service
	analyticsService         *service_analytics.Service
	viewService              *service_view.Service
	snapshotter              *service_snapshots.Service
	commentService           *service_comments.Service
	activityService          *service_activity.Service
	codebaseService          *service_codebase.Service
	workspaceStatusesService *service_workspace_statuses.Service

	activitySender   sender.ActivitySender
	snapshotterQueue worker_snapshots.Queue
	eventsSender     events.EventSender
	eventsPublisher  *eventsv2.Publisher
	executorProvider executor.Provider
	buildQueue       *workers_ci.BuildQueue
}

func New(
	logger *zap.Logger,

	workspaceWriter db_workspaces.WorkspaceWriter,

	usersService service_users.Service,
	workspaceService *service_workspaces.Service,
	changeService *service_changes.Service,
	analyticsService *service_analytics.Service,
	viewService *service_view.Service,
	snapshotter *service_snapshots.Service,
	commentService *service_comments.Service,
	activityService *service_activity.Service,
	codebaseService *service_codebase.Service,
	workspaceStatusesService *service_workspace_statuses.Service,

	activitySender sender.ActivitySender,
	snapshotterQueue worker_snapshots.Queue,
	eventsSender events.EventSender,
	eventsPublisher *eventsv2.Publisher,
	executorProvider executor.Provider,
	buildQueue *workers_ci.BuildQueue,
) *Service {
	return &Service{
		logger: logger,

		workspaceWriter: workspaceWriter,

		usersService:             usersService,
		workspaceService:         workspaceService,
		changeService:            changeService,
		analyticsService:         analyticsService,
		viewService:              viewService,
		snapshotter:              snapshotter,
		commentService:           commentService,
		activityService:          activityService,
		codebaseService:          codebaseService,
		workspaceStatusesService: workspaceStatusesService,

		activitySender:   activitySender,
		snapshotterQueue: snapshotterQueue,
		eventsSender:     eventsSender,
		eventsPublisher:  eventsPublisher,
		executorProvider: executorProvider,
		buildQueue:       buildQueue,
	}
}

func (s *Service) LandChange(ctx context.Context, ws *workspaces.Workspace, diffOpts ...vcs.DiffOption) (*changes.Change, error) {
	user, err := s.usersService.GetByID(ctx, ws.UserID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	// check if the workspace is allowed to be landed
	cb, err := s.codebaseService.GetByID(ctx, ws.CodebaseID)
	if err != nil {
		return nil, fmt.Errorf("failed to get codebase: %w", err)
	}

	// make sure that all statuses are healthy and not stale
	if cb.RequireHealthyStatus {
		healthy, err := s.workspaceStatusesService.HealthyStatus(ctx, ws)
		switch {
		case err != nil:
			return nil, fmt.Errorf("failed to get workspace status: %w", err)
		case !healthy:
			return nil, ErrNotAllowedUnhealthyWorkspace
		}
	}

	gitCommitMessage := message.CommitMessage(ws.DraftDescription)

	signature := git.Signature{
		Name:  user.Name,
		Email: user.Email,
		When:  time.Now(),
	}

	var change *changes.Change
	creteAndLand := func(viewRepo vcs.RepoWriter) error {
		createdCommitID, fromViewPushFunc, err := s.changeService.CreateAndLandFromView(
			context.Background(),
			viewRepo,
			ws.CodebaseID,
			ws.ID,
			gitCommitMessage,
			signature,
			diffOpts...,
		)
		if err != nil {
			return fmt.Errorf("failed to create and land from view: %w", err)
		}

		parents, err := viewRepo.GetCommitParents(createdCommitID)
		if err != nil {
			return fmt.Errorf("failed get parents of new commit: %w", err)
		}
		if len(parents) != 1 {
			return fmt.Errorf("commit has an unexpected number of parents n=%d", len(parents))
		}

		change, err = s.changeService.CreateWithCommitAsParent(ctx, ws, createdCommitID, parents[0])
		if err != nil {
			return fmt.Errorf("failed to create change: %w", err)
		}

		if err := fromViewPushFunc(viewRepo); err != nil {
			return fmt.Errorf("failed to push the landed result: %w", err)
		}
		return nil
	}

	if ws.ViewID != nil {
		if err := s.executorProvider.New().
			Write(creteAndLand).
			ExecView(ws.CodebaseID, *ws.ViewID, "landChangeCreateAndLandFromView"); err != nil {
			return nil, fmt.Errorf("failed to share from view: %w", err)
		}
	} else {
		if ws.LatestSnapshotID == nil {
			return nil, fmt.Errorf("the workspace has no snapshot")
		}
		snapshot, err := s.snapshotter.GetByID(ctx, *ws.LatestSnapshotID)
		if err != nil {
			return nil, fmt.Errorf("failed to get snapshot: %w", err)
		}
		if err := s.executorProvider.New().
			Write(func(writer vcs.RepoWriter) error {
				return writer.CreateBranchTrackingUpstream(ws.ID)
			}).
			Write(vcs_view.CheckoutSnapshot(snapshot)).
			Write(creteAndLand).
			ExecTemporaryView(ws.CodebaseID, "landChangeCreateAndLandFromSnapshot"); err != nil {
			return nil, fmt.Errorf("failed to create and land from snaphsot: %w", err)
		}
		ws.SetSnapshot(nil)
	}

	s.analyticsService.Capture(ctx, "create change",
		analytics.CodebaseID(ws.CodebaseID),
		analytics.Property("workspace_id", ws.ID),
		analytics.Property("change_id", change.ID),
	)

	if ws.ViewID != nil {
		if err := s.snapshotterQueue.Enqueue(ctx, ws.CodebaseID, *ws.ViewID, ws.ID, ws.UserID, snapshots.ActionChangeLand); err != nil {
			return nil, fmt.Errorf("failed to enqueue snapshot: %w", err)
		}

		vw, err := s.viewService.GetByID(ctx, *ws.ViewID)
		if err != nil {
			return nil, fmt.Errorf("could not get view: %w", err)
		}

		if err := s.eventsPublisher.ViewUpdated(ctx, eventsv2.Codebase(ws.CodebaseID), vw); err != nil {
			return nil, fmt.Errorf("failed to send view updated event: %w", err)
		}
	}

	// Update workspace
	now := time.Now()
	if err := s.workspaceWriter.UpdateFields(ctx, ws.ID,
		db_workspaces.SetHeadChangeID(nil), // TODO: Set this directly
		db_workspaces.SetHeadChangeComputed(false),
		db_workspaces.SetUpdatedAt(&now),
		db_workspaces.SetDraftDescription(""),
		db_workspaces.SetChangeID(&change.ID),
		db_workspaces.SetLastLandedAt(&now),
	); err != nil {
		return nil, fmt.Errorf("failed to update workspace: %w", err)
	}

	// Send event that the workspace has been updated
	if err := s.eventsSender.Workspace(ws.ID, events.WorkspaceUpdated, ws.ID); err != nil {
		s.logger.Error("failed to send workspace event", zap.Error(err))
	}

	// Clear 'up to date' cache for all workspaces
	if err := s.workspaceWriter.UnsetUpToDateWithTrunkForAllInCodebase(ws.CodebaseID); err != nil {
		return nil, fmt.Errorf("failed to unset up_to_date_with_trunk: %w", err)
	}

	s.analyticsService.Capture(ctx, "landed changes",
		analytics.CodebaseID(ws.CodebaseID),
		analytics.Property("workspace_id", ws.ID),
		analytics.Property("change_id", change.ID),
	)

	if err := s.commentService.MoveCommentsFromWorkspaceToChange(ctx, ws.ID, change.ID); err != nil {
		return nil, fmt.Errorf("failed to move comments from workspace to change: %w", err)
	}

	// Create activity
	if err := s.activitySender.Codebase(ctx, ws.CodebaseID, ws.ID, ws.UserID, activity.TypeCreatedChange, string(change.ID)); err != nil {
		return nil, fmt.Errorf("failed to create workspace activity: %w", err)
	}

	// Make activity list available for the change
	if err := s.activityService.SetChange(ctx, ws.ID, change.ID); err != nil {
		return nil, fmt.Errorf("failed to set change activity: %w", err)
	}

	// Update codebase cache
	if err := s.changeService.SetAsHeadChange(change); err != nil {
		return nil, fmt.Errorf("failed to set as head change: %w", err)
	}

	// Send events that the codebase has been updated
	if err := s.eventsSender.Codebase(ws.CodebaseID, events.CodebaseUpdated, ws.CodebaseID.String()); err != nil {
		s.logger.Error("failed to send codebase event", zap.Error(err))
	}

	if err := s.eventsSender.Workspace(ws.ID, events.WorkspaceUpdatedSnapshot, ws.ID); err != nil {
		s.logger.Error("failed to send workspace event", zap.Error(err))
	}

	if err := s.buildQueue.EnqueueChange(ctx, change); err != nil {
		s.logger.Error("failed to enqueue change", zap.Error(err))
	}

	if err := s.workspaceService.ArchiveWithChange(ctx, ws, change); err != nil {
		return nil, fmt.Errorf("failed to archive workspace: %w", err)
	}

	return change, nil
}
