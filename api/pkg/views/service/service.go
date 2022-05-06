package service

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"getsturdy.com/api/pkg/analytics"
	service_analytics "getsturdy.com/api/pkg/analytics/service"
	events "getsturdy.com/api/pkg/events/v2"
	"getsturdy.com/api/pkg/snapshots"
	db_snapshots "getsturdy.com/api/pkg/snapshots/db"
	service_snapshots "getsturdy.com/api/pkg/snapshots/service"
	"getsturdy.com/api/pkg/users"
	"getsturdy.com/api/pkg/views"
	"getsturdy.com/api/pkg/views/db"
	vcs_view "getsturdy.com/api/pkg/views/vcs"
	"getsturdy.com/api/pkg/workspaces"
	db_workspaces "getsturdy.com/api/pkg/workspaces/db"
	"getsturdy.com/api/vcs"
	"getsturdy.com/api/vcs/executor"

	"github.com/google/uuid"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

type Service struct {
	logger           *zap.Logger
	viewRepo         db.Repository
	workspaceReader  db_workspaces.WorkspaceReader
	gitSnapshotter   *service_snapshots.Service
	snapshotRepo     db_snapshots.Repository
	workspaceWriter  db_workspaces.WorkspaceWriter
	executorProvider executor.Provider
	eventSender      *events.Publisher
	analyticsService *service_analytics.Service
}

func New(
	logger *zap.Logger,
	viewRepo db.Repository,
	workspaceReader db_workspaces.WorkspaceReader,
	gitSnapshotter *service_snapshots.Service,
	snapshotRepo db_snapshots.Repository,
	workspaceWriter db_workspaces.WorkspaceWriter,
	executorProvider executor.Provider,
	eventSender *events.Publisher,
	analyticsService *service_analytics.Service,
) *Service {
	return &Service{
		logger:           logger.Named("views_service"),
		viewRepo:         viewRepo,
		workspaceReader:  workspaceReader,
		gitSnapshotter:   gitSnapshotter,
		snapshotRepo:     snapshotRepo,
		workspaceWriter:  workspaceWriter,
		executorProvider: executorProvider,
		eventSender:      eventSender,
		analyticsService: analyticsService,
	}
}

var ErrRebasing = errors.New("is rebasing")

func (s *Service) OpenWorkspace(ctx context.Context, view *views.View, ws *workspaces.Workspace) error {
	if ws.ArchivedAt != nil {
		return fmt.Errorf("the workspace is archived")
	}
	if view.CodebaseID != ws.CodebaseID {
		return fmt.Errorf("the view and workspace does not belong to the same codebase")
	}
	if view.UserID != ws.UserID {
		return fmt.Errorf("the view and workspace does not belong to the same user")
	}

	// If the view that we're opening this workspace on has another workspace currently open, snapshot and save the changes
	currentWorkspaceOnView, err := s.workspaceReader.GetByViewID(view.ID, true)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return fmt.Errorf("could not find previous workspace on view: %w", err)
	} else if err == nil {
		_, err := s.gitSnapshotter.Snapshot(ctx, currentWorkspaceOnView.CodebaseID, currentWorkspaceOnView.ID,
			snapshots.ActionPreCheckoutOtherWorkspace, service_snapshots.WithOnView(*currentWorkspaceOnView.ViewID), service_snapshots.WithMarkAsLatestInWorkspace())
		if errors.Is(err, service_snapshots.ErrCantSnapshotRebasing) {
			return ErrRebasing
		}
		if err != nil {
			return fmt.Errorf("failed to snapshot: %w", err)
		}

		currentWorkspaceOnView.ViewID = nil
		if err := s.workspaceWriter.UpdateFields(ctx, currentWorkspaceOnView.ID, db_workspaces.SetViewID(nil)); err != nil {
			return fmt.Errorf("failed to finalize previous workspace on view: %w", err)
		}
	}

	// If the workspace that we're opening on the view, currently is opened somewhere else, snapshot it and save the changes
	// The snapshot will also be used to "move" the contents to the new view
	if ws.ViewID != nil {
		_, err := s.gitSnapshotter.Snapshot(ctx, ws.CodebaseID, ws.ID, snapshots.ActionPreCheckoutOtherView,
			service_snapshots.WithOnView(*ws.ViewID), service_snapshots.WithMarkAsLatestInWorkspace())
		if err != nil {
			return fmt.Errorf("failed to snapshot: %w", err)
		}

		// TODO: unset view.workspace_id?
	}

	if err := s.executorProvider.New().
		Write(vcs_view.CheckoutBranch(ws.ID)).
		Write(func(repo vcs.RepoWriter) error {
			// Restore snapshot
			if ws.LatestSnapshotID != nil {
				snapshot, err := s.snapshotRepo.Get(*ws.LatestSnapshotID)
				if err != nil {
					return fmt.Errorf("failed to get snapshot: %w", err)
				}
				if err := s.gitSnapshotter.Restore(snapshot, repo); err != nil {
					return fmt.Errorf("failed to restore snapshot: %w", err)
				}
			} else {
				// vcs2.Restore does this as well, make sure to smudge files in this scenario as well
				if err := repo.LargeFilesPull(); err != nil {
					s.logger.Warn("failed to pull large files", zap.Error(err))
				}
				if err := repo.CleanStaged(); err != nil {
					return fmt.Errorf("failed to clean index: %w", err)
				}
			}
			return nil
		}).ExecView(ws.CodebaseID, view.ID, "OpenWorkspaceOnView"); err != nil {
		return fmt.Errorf("failed to open workspace on view: %w", err)
	}

	// Update workspace object
	ws.ViewID = &view.ID
	if err := s.workspaceWriter.UpdateFields(ctx, ws.ID, db_workspaces.SetViewID(&view.ID)); err != nil {
		return fmt.Errorf("failed to update workspace: %w", err)
	}

	// Update view object
	view.WorkspaceID = ws.ID
	if err := s.viewRepo.Update(view); err != nil {
		return fmt.Errorf("failed to update view obj: %w", err)
	}

	if err := s.eventSender.ViewUpdated(ctx, events.Workspace(ws.ID), view); err != nil {
		s.logger.Error("failed to send view updated event: %w", zap.Error(err))
		// do not fail
	}

	return nil
}

func (s *Service) GetByID(_ context.Context, id string) (*views.View, error) {
	return s.viewRepo.Get(id)
}

func (s *Service) Create(ctx context.Context, userID users.ID, workspace *workspaces.Workspace, mountPath, mountHostname *string) (*views.View, error) {
	t := time.Now()
	v := views.View{
		ID:            uuid.New().String(),
		UserID:        userID,
		CodebaseID:    workspace.CodebaseID,
		WorkspaceID:   workspace.ID,
		MountPath:     mountPath,     // It's optional
		MountHostname: mountHostname, // It's optional
		CreatedAt:     &t,
	}

	if err := s.viewRepo.Create(v); err != nil {
		return nil, fmt.Errorf("failed to create view: %w", err)
	}

	if err := s.executorProvider.New().
		AllowRebasingState(). // allowed because the view does not exist yet
		Schedule(vcs_view.Create(workspace.CodebaseID, workspace.ID, v.ID)).
		ExecView(workspace.CodebaseID, v.ID, "createView"); err != nil {
		return nil, fmt.Errorf("failed to create view: %w", err)
	}

	// Use workspace on view
	if err := s.OpenWorkspace(ctx, &v, workspace); err != nil {
		return nil, fmt.Errorf("failed to open workspace on view: %w", err)
	}

	s.analyticsService.Capture(ctx, "create view",
		analytics.CodebaseID(workspace.CodebaseID),
		analytics.Property("workspace_id", workspace.ID),
		analytics.Property("view_id", v.ID),
		analytics.Property("mount_path", v.MountPath),
		analytics.Property("mount_hostname", v.MountHostname),
	)

	return &v, nil
}
