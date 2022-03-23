package service

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	"getsturdy.com/api/pkg/activity"
	"getsturdy.com/api/pkg/activity/sender"
	service_activity "getsturdy.com/api/pkg/activity/service"
	"getsturdy.com/api/pkg/analytics"
	service_analytics "getsturdy.com/api/pkg/analytics/service"
	"getsturdy.com/api/pkg/changes"
	"getsturdy.com/api/pkg/changes/message"
	service_change "getsturdy.com/api/pkg/changes/service"
	change_vcs "getsturdy.com/api/pkg/changes/vcs"
	workers_ci "getsturdy.com/api/pkg/ci/workers"
	"getsturdy.com/api/pkg/codebases"
	service_comments "getsturdy.com/api/pkg/comments/service"
	"getsturdy.com/api/pkg/events"
	eventsv2 "getsturdy.com/api/pkg/events/v2"
	db_review "getsturdy.com/api/pkg/review/db"
	"getsturdy.com/api/pkg/snapshots"
	"getsturdy.com/api/pkg/snapshots/snapshotter"
	worker_snapshots "getsturdy.com/api/pkg/snapshots/worker"
	"getsturdy.com/api/pkg/unidiff"
	"getsturdy.com/api/pkg/unidiff/lfs"
	"getsturdy.com/api/pkg/users"
	user_db "getsturdy.com/api/pkg/users/db"
	service_view "getsturdy.com/api/pkg/view/service"
	vcs_view "getsturdy.com/api/pkg/view/vcs"
	"getsturdy.com/api/pkg/workspaces"
	"getsturdy.com/api/pkg/workspaces/db"
	vcs_workspace "getsturdy.com/api/pkg/workspaces/vcs"
	"getsturdy.com/api/vcs"
	"getsturdy.com/api/vcs/executor"

	"github.com/google/uuid"
	git "github.com/libgit2/git2go/v33"
	"go.uber.org/zap"
)

type CreateWorkspaceRequest struct {
	UserID           users.ID
	CodebaseID       codebases.ID
	Name             string
	DraftDescription string

	BaseChangeID *changes.ID
	Revert       bool
}

type Service interface {
	Create(context.Context, CreateWorkspaceRequest) (*workspaces.Workspace, error)
	CreateFromWorkspace(ctx context.Context, from *workspaces.Workspace, userID users.ID, name string) (*workspaces.Workspace, error)
	GetByID(context.Context, string) (*workspaces.Workspace, error)
	GetByViewID(context.Context, string) (*workspaces.Workspace, error)
	LandChange(ctx context.Context, ws *workspaces.Workspace, patchIDs []string, diffOptions ...vcs.DiffOption) (*changes.Change, error)
	CreateWelcomeWorkspace(ctx context.Context, codebaseID codebases.ID, userID users.ID, codebaseName string) error
	Diffs(context.Context, string, ...DiffsOption) ([]unidiff.FileDiff, bool, error)
	CopyPatches(ctx context.Context, src, dist *workspaces.Workspace, opts ...CopyPatchesOption) error
	RemovePatches(context.Context, *unidiff.Allower, *workspaces.Workspace, ...string) error
	HasConflicts(context.Context, *workspaces.Workspace) (bool, error)
	Archive(context.Context, *workspaces.Workspace) error
	ArchiveWithChange(context.Context, *workspaces.Workspace, *changes.Change) error
	Unarchive(context.Context, *workspaces.Workspace) error
	HeadChange(ctx context.Context, ws *workspaces.Workspace) (*changes.Change, error)
	Push(ctx context.Context, user *users.User, ws *workspaces.Workspace) error
	LandOnSturdyAndPushTracked(ctx context.Context, ws *workspaces.Workspace) error
}

type WorkspaceService struct {
	logger           *zap.Logger
	analyticsService *service_analytics.Service

	workspaceWriter db.WorkspaceWriter
	workspaceReader db.WorkspaceReader

	userRepo   user_db.Repository
	reviewRepo db_review.ReviewRepository

	commentService  *service_comments.Service
	changeService   *service_change.Service
	activityService *service_activity.Service
	viewService     *service_view.Service

	activitySender   sender.ActivitySender
	eventsSender     events.EventSender
	eventsSenderV2   *eventsv2.Publisher
	snapshotterQueue worker_snapshots.Queue
	executorProvider executor.Provider
	snap             snapshotter.Snapshotter
	buildQueue       *workers_ci.BuildQueue
}

func New(
	logger *zap.Logger,
	analyticsService *service_analytics.Service,

	workspaceWriter db.WorkspaceWriter,
	workspaceReader db.WorkspaceReader,

	userRepo user_db.Repository,
	reviewRepo db_review.ReviewRepository,

	commentsService *service_comments.Service,
	changeService *service_change.Service,
	activityService *service_activity.Service,
	viewService *service_view.Service,

	activitySender sender.ActivitySender,
	executorProvider executor.Provider,
	eventsSender events.EventSender,
	eventsSenderV2 *eventsv2.Publisher,
	snapshotterQueue worker_snapshots.Queue,
	snap snapshotter.Snapshotter,
	buildQueue *workers_ci.BuildQueue,
) *WorkspaceService {
	return &WorkspaceService{
		logger:           logger,
		analyticsService: analyticsService,

		workspaceWriter: workspaceWriter,
		workspaceReader: workspaceReader,

		userRepo:   userRepo,
		reviewRepo: reviewRepo,

		commentService:  commentsService,
		changeService:   changeService,
		activityService: activityService,
		viewService:     viewService,

		activitySender:   activitySender,
		executorProvider: executorProvider,
		eventsSender:     eventsSender,
		eventsSenderV2:   eventsSenderV2,
		snapshotterQueue: snapshotterQueue,
		snap:             snap,
		buildQueue:       buildQueue,
	}
}

type DiffsOptions struct {
	Allower        *unidiff.Allower
	VCSDiffOptions []vcs.DiffOption
}

type DiffsOption func(*DiffsOptions)

func WithVCSDiffOptions(options ...vcs.DiffOption) DiffsOption {
	return func(diffsOptions *DiffsOptions) {
		diffsOptions.VCSDiffOptions = append(diffsOptions.VCSDiffOptions, options...)
	}
}

func WithAllower(allower *unidiff.Allower) DiffsOption {
	return func(options *DiffsOptions) {
		options.Allower = allower
	}
}

func getDiffOptions(opts ...DiffsOption) *DiffsOptions {
	options := &DiffsOptions{}
	for _, opt := range opts {
		opt(options)
	}
	return options
}

func (s *WorkspaceService) GetByViewID(ctx context.Context, viewID string) (*workspaces.Workspace, error) {
	return s.workspaceReader.GetByViewID(viewID, true)
}

func (s *WorkspaceService) Diffs(ctx context.Context, workspaceID string, oo ...DiffsOption) ([]unidiff.FileDiff, bool, error) {
	ws, err := s.GetByID(ctx, workspaceID)
	if err != nil {
		return nil, false, fmt.Errorf("failed to find workspace: %w", err)
	}
	options := getDiffOptions(oo...)
	if ws.ViewID == nil {
		diffs, err := s.diffsFromSnapshot(ctx, ws, options)
		return diffs, false, err
	}

	return s.diffsFromView(ctx, ws, options)
}

func (s *WorkspaceService) diffsFromSnapshot(ctx context.Context, ws *workspaces.Workspace, options *DiffsOptions) ([]unidiff.FileDiff, error) {
	if ws.LatestSnapshotID == nil {
		return nil, nil
	}

	snapshotOptions := []snapshotter.DiffsOption{}
	if options.Allower != nil {
		snapshotOptions = append(snapshotOptions, snapshotter.WithAllower(options.Allower))
	}

	return s.snap.Diffs(ctx, *ws.LatestSnapshotID, snapshotOptions...)
}

func (s *WorkspaceService) diffsFromView(ctx context.Context, ws *workspaces.Workspace, options *DiffsOptions) ([]unidiff.FileDiff, bool, error) {
	var diffs []unidiff.FileDiff

	isRebasing := false
	if err := s.executorProvider.New().
		AssertBranchName(ws.ID).
		AllowRebasingState(). // allowed to generate diffs even if conflicting
		Read(func(repo vcs.RepoReader) error {
			isRebasing = repo.IsRebasing()

			gitDiffs, err := repo.Diffs(options.VCSDiffOptions...)
			if err != nil {
				return fmt.Errorf("failed to get git repo diffs: %w", err)
			}
			defer gitDiffs.Free()

			filter, err := lfs.NewIgnoreLfsSmudgedFilter(repo)
			if err != nil {
				return fmt.Errorf("could not smudge lfs files: %w", err)
			}

			differ := unidiff.NewUnidiff(unidiff.NewGitPatchReader(gitDiffs), s.logger).
				WithExpandedHunks().
				WithFilterFunc(filter)

			if options.Allower != nil {
				differ = differ.WithAllower(options.Allower)
			}

			hunkifiedDiff, err := differ.Decorate()
			if err != nil {
				return fmt.Errorf("could not decorate view diffs: %w", err)
			}

			diffs = hunkifiedDiff
			return nil
		}).ExecView(ws.CodebaseID, *ws.ViewID, "workspaceViewDiffs"); err != nil {
		return nil, false, fmt.Errorf("failed to get diffs from view: %w", err)
	}
	return diffs, isRebasing, nil
}

func (s *WorkspaceService) GetByID(ctx context.Context, id string) (*workspaces.Workspace, error) {
	return s.workspaceReader.Get(id)
}

type CopyPatchesOptions struct {
	PatchIDs *[]string
}

type CopyPatchesOption func(*CopyPatchesOptions)

func WithPatchIDs(patchIDs []string) CopyPatchesOption {
	return func(options *CopyPatchesOptions) {
		if options.PatchIDs == nil {
			options.PatchIDs = new([]string)
		}
		*options.PatchIDs = append(*options.PatchIDs, patchIDs...)
	}
}

func getCopyPatchOptions(oo ...CopyPatchesOption) *CopyPatchesOptions {
	options := &CopyPatchesOptions{}
	for _, o := range oo {
		o(options)
	}
	return options
}

func (s *WorkspaceService) CopyPatches(ctx context.Context, dist, src *workspaces.Workspace, opts ...CopyPatchesOption) error {
	if src.CodebaseID != dist.CodebaseID {
		return fmt.Errorf("source and destination codebases must be the same")
	}

	if dist.ViewID != nil {
		// TODO
		return fmt.Errorf("copying to active workspace is not supported")
	}

	options := getCopyPatchOptions(opts...)
	fields := []db.UpdateOption{}
	if src.ViewID != nil {
		// if workspace has a view, snapshot changes from it
		snapshotterOptions := []snapshotter.SnapshotOption{snapshotter.WithOnView(*src.ViewID)}
		if options.PatchIDs != nil {
			snapshotterOptions = append(snapshotterOptions, snapshotter.WithPatchIDsFilter(*options.PatchIDs))
		}
		snapshot, err := s.snap.Snapshot(src.CodebaseID, src.ID, snapshots.ActionWorkspaceExtract, snapshotterOptions...)
		if err != nil {
			return fmt.Errorf("failed to create snapshot: %w", err)
		}
		fields = append(fields, db.SetLatestSnapshotID(&snapshot.ID), db.SetDiffsCount(snapshot.DiffsCount))
	} else if options.PatchIDs != nil {
		// if workspace doesn't have a view, copy patches from it's latest snapshot
		if src.LatestSnapshotID == nil {
			return fmt.Errorf("source workspace doesn't have a snapshot")
		}
		copyOptions := []snapshotter.CopyOption{}
		if options.PatchIDs != nil {
			copyOptions = append(copyOptions, snapshotter.CopyWithPatchIDs(*options.PatchIDs))
		}
		snapshot, err := s.snap.Copy(ctx, *src.LatestSnapshotID, copyOptions...)
		if err != nil {
			return fmt.Errorf("failed to copy snapshot: %w", err)
		}
		fields = append(fields, db.SetLatestSnapshotID(&snapshot.ID), db.SetDiffsCount(snapshot.DiffsCount))
	} else {
		// if we don't need to copy patches, re-use the existing snapshot
		fields = append(fields, db.SetLatestSnapshotID(src.LatestSnapshotID), db.SetDiffsCount(src.DiffsCount))
	}

	if err := s.workspaceWriter.UpdateFields(ctx, dist.ID, fields...); err != nil {
		return fmt.Errorf("failed to update workspace: %w", err)
	}

	return nil
}

func (s *WorkspaceService) CreateFromWorkspace(ctx context.Context, from *workspaces.Workspace, userID users.ID, name string) (*workspaces.Workspace, error) {

	var baseChangeID *changes.ID
	fromBaseChange, err := s.HeadChange(ctx, from)
	switch {
	case errors.Is(err, ErrNotFound):
	// head change not found (this repo does not have any changes, or this workspace is based on the root)
	// do nothing
	case err != nil:
		return nil, fmt.Errorf("failed to get head change: %w", err)
	default:
		baseChangeID = &fromBaseChange.ID
	}

	createRequest := CreateWorkspaceRequest{
		UserID:       userID,
		CodebaseID:   from.CodebaseID,
		Name:         name,
		BaseChangeID: baseChangeID,
	}

	newWorkspace, err := s.Create(ctx, createRequest)
	if err != nil {
		return nil, fmt.Errorf("faliled to create a workspace: %w", err)
	}

	return newWorkspace, nil
}

func (s *WorkspaceService) Create(ctx context.Context, req CreateWorkspaceRequest) (*workspaces.Workspace, error) {
	t := time.Now()
	var zero int32 = 0
	ws := workspaces.Workspace{
		ID:               uuid.New().String(),
		UserID:           req.UserID,
		CodebaseID:       req.CodebaseID,
		CreatedAt:        &t,
		DraftDescription: req.DraftDescription,
		DiffsCount:       &zero,
	}

	if len(req.Name) > 0 {
		ws.Name = &req.Name
	} else {
		n := "Untitled draft"
		ws.Name = &n
	}

	var baseCommitSha string
	var baseCommitParentSha *string
	if req.BaseChangeID != nil {
		ch, err := s.changeService.GetChangeByID(ctx, *req.BaseChangeID)
		if err != nil {
			return nil, fmt.Errorf("could not get change by id: %w", err)
		}
		if ch.CodebaseID != ws.CodebaseID {
			return nil, fmt.Errorf("change does not belong to this codebase")
		}
		if ch.CommitID == nil {
			return nil, fmt.Errorf("the change does not have a commit")
		}

		baseCommitSha = *ch.CommitID

		// If the change has a parent, calculate the diffs between the change and it's parent
		// otherwise use the diff between the change and the root of the repo
		if ch.ParentChangeID != nil {
			parentChange, err := s.changeService.GetChangeByID(ctx, *ch.ParentChangeID)
			if err != nil {
				return nil, fmt.Errorf("could not get parent change by id: %w", err)
			}
			if parentChange.CommitID == nil {
				return nil, fmt.Errorf("the change parent does not have a commit")
			}
			baseCommitParentSha = parentChange.CommitID
		}
	}

	if err := s.executorProvider.New().GitWrite(func(repo vcs.RepoGitWriter) error {
		// Ensure codebase status
		if err := EnsureCodebaseStatus(repo); err != nil {
			return err
		}

		if req.BaseChangeID != nil && baseCommitSha != "" {
			// Create workspace at the change that we want to revert
			if err := vcs_workspace.CreateOnCommitID(repo, ws.ID, baseCommitSha); err != nil {
				return fmt.Errorf("failed to create workspace at change: %w", err)
			}
		} else {
			// Create workspace at current trunk
			if err := vcs_workspace.Create(repo, ws.ID); err != nil {
				return fmt.Errorf("failed to create workspace: %w", err)
			}
		}
		return nil
	}).ExecTrunk(req.CodebaseID, "createWorkspace"); err != nil {
		return nil, err
	}

	if err := s.workspaceWriter.Create(ws); err != nil {
		return nil, fmt.Errorf("failed to write workspace to db: %w", err)
	}

	// Add the reverted changes to a snapshot
	if req.BaseChangeID != nil && baseCommitSha != "" && req.Revert {
		snapshot, err := s.snap.Snapshot(
			ws.CodebaseID,
			ws.ID,
			snapshots.ActionChangeReverted,
			snapshotter.WithOnTemporaryView(),
			snapshotter.WithRevertDiff(baseCommitSha, baseCommitParentSha),
		)
		if err != nil {
			return nil, fmt.Errorf("failed to create snapshot for revert: %w", err)
		}
		ws.SetSnapshot(snapshot)
		if err := s.workspaceWriter.UpdateFields(ctx, ws.ID, db.SetLatestSnapshotID(ws.LatestSnapshotID), db.SetDiffsCount(ws.DiffsCount)); err != nil {
			return nil, fmt.Errorf("failed to update workspace: %w", err)
		}
	}

	s.analyticsService.Capture(ctx, "create workspace",
		analytics.CodebaseID(req.CodebaseID),
		analytics.Property("id", ws.ID),
		analytics.Property("at_existing_change", req.BaseChangeID != nil),
		analytics.Property("name", ws.Name),
	)

	return &ws, nil
}

var ErrNotFound = errors.New("not found")

func (s *WorkspaceService) HeadChange(ctx context.Context, ws *workspaces.Workspace) (*changes.Change, error) {
	if ws.HeadChangeComputed {
		if ws.HeadChangeID == nil {
			return nil, ErrNotFound
		}
		ch, err := s.changeService.GetChangeByID(ctx, *ws.HeadChangeID)
		if err != nil {
			return nil, err
		}
		return ch, nil
	}

	// Compute!
	var headCommitID string

	err := s.executorProvider.New().GitRead(func(repo vcs.RepoGitReader) error {
		var err error
		headCommitID, err = repo.BranchCommitID(ws.ID)
		if err != nil {
			return fmt.Errorf("could not get head commit from git: %w", err)
		}
		return nil
	}).ExecTrunk(ws.CodebaseID, "workspaceHeadChange")
	if err != nil {
		return nil, err
	}
	var newHeadChangeID *changes.ID

	ch, err := s.changeService.GetByCommitAndCodebase(ctx, headCommitID, ws.CodebaseID)
	switch {
	case errors.Is(err, sql.ErrNoRows), errors.Is(err, service_change.ErrNotFound):
		// change not found (could be the root commit, etc), hide it
		newHeadChangeID = nil
	case err != nil:
		return nil, fmt.Errorf("could not get change by commit: %w", err)
	default:
		newHeadChangeID = &ch.ID
	}

	// Save updated cache
	if err := s.workspaceWriter.UpdateFields(ctx, ws.ID,
		db.SetHeadChangeID(newHeadChangeID),
		db.SetHeadChangeComputed(true),
	); err != nil {
		return nil, err
	}

	// Also update the cached version of the workspace that we have in memory
	ws.HeadChangeComputed = true
	ws.HeadChangeID = newHeadChangeID

	s.logger.Info("recalculated head change", zap.String("workspace_id", ws.ID), zap.Stringer("head", ws.HeadChangeID))

	if ch == nil {
		return nil, ErrNotFound
	}

	return ch, nil
}

func (s *WorkspaceService) LandChange(ctx context.Context, ws *workspaces.Workspace, patchIDs []string, diffOpts ...vcs.DiffOption) (*changes.Change, error) {
	user, err := s.userRepo.Get(ws.UserID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	gitCommitMessage := message.CommitMessage(ws.DraftDescription)

	signature := git.Signature{
		Name:  user.Name,
		Email: user.Email,
		When:  time.Now(),
	}

	var change *changes.Change
	creteAndLand := func(viewRepo vcs.RepoWriter) error {
		createdCommitID, fromViewPushFunc, err := change_vcs.CreateAndLandFromView(
			viewRepo,
			s.logger,
			ws.CodebaseID,
			ws.ID,
			patchIDs,
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
		snapshot, err := s.snap.GetByID(ctx, *ws.LatestSnapshotID)
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
		if err := s.snapshotterQueue.Enqueue(ctx, ws.CodebaseID, *ws.ViewID, ws.ID, snapshots.ActionChangeLand); err != nil {
			return nil, fmt.Errorf("failed to enqueue snapshot: %w", err)
		}

		vw, err := s.viewService.GetByID(ctx, *ws.ViewID)
		if err != nil {
			return nil, fmt.Errorf("could not get view: %w", err)
		}

		if err := s.eventsSenderV2.ViewUpdated(ctx, eventsv2.Codebase(ws.CodebaseID), vw); err != nil {
			return nil, fmt.Errorf("failed to send view updated event: %w", err)
		}
	}

	// Update workspace
	now := time.Now()
	if err := s.workspaceWriter.UpdateFields(ctx, ws.ID,
		db.SetHeadChangeID(nil), // TODO: Set this directly
		db.SetHeadChangeComputed(false),
		db.SetUpdatedAt(&now),
		db.SetDraftDescription(""),
		db.SetChangeID(&change.ID),
		db.SetLastLandedAt(&now),
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

	if err := s.Archive(ctx, ws); err != nil {
		return nil, fmt.Errorf("failed to archive workspace: %w", err)
	}

	return change, nil
}

func EnsureCodebaseStatus(repo vcs.RepoGitWriter) error {
	// Make sure that a root commit exists
	// This is the first time a root commit is _needed_ (so that we can create a branch),
	// and we don't want to do it earlier (such as on clone from GitHub).
	// If there is no head / root commit, create one

	if _, err := repo.HeadCommit(); err != nil {
		if err := repo.CreateRootCommit(); err != nil {
			return err
		}
	}

	// If sturdytrunk is not the default branch, create it
	defaultBranch, err := repo.GetDefaultBranch()
	if err != nil {
		return err
	}
	if defaultBranch != "refs/heads/sturdytrunk" {
		if err := repo.CreateAndSetDefaultBranch("sturdytrunk"); err != nil {
			return err
		}
	}

	return nil
}

const readMeTemplate = `# __CODEBASE__NAME__
`

const draftDescriptionTemplate = `<h3>Adding a README to __CODEBASE__NAME__</h3>
<ul>
	<li><p>This is a workspace - it's where you're <strong>coding</strong>, and can give and take <strong>feedback</strong> from your team</p></li>
	<li><p><strong>Share</strong> this draft to land the changes to the README on the trunk, and to make the file available to all collaborators</p></li>
</ul>

<p>Happy hacking!</p>
`

func (svc *WorkspaceService) CreateWelcomeWorkspace(ctx context.Context, codebaseID codebases.ID, userID users.ID, codebaseName string) error {
	readMeContents := strings.ReplaceAll(readMeTemplate, "__CODEBASE__NAME__", codebaseName)
	draftDescriptionContents := strings.ReplaceAll(draftDescriptionTemplate, "__CODEBASE__NAME__", codebaseName)

	ws, err := svc.Create(ctx, CreateWorkspaceRequest{
		CodebaseID:       codebaseID,
		UserID:           userID,
		Name:             "Add README",
		DraftDescription: draftDescriptionContents,
	})
	if err != nil {
		return fmt.Errorf("failed to create workspace: %w", err)
	}

	cb := func(repo vcs.RepoReaderGitWriter) error {
		branchName := "welcome-" + uuid.NewString()

		commitID, err := repo.CreateCommitWithFiles([]vcs.FileContents{
			{Path: "README.md", Contents: []byte(readMeContents)},
		}, branchName)
		if err != nil {
			return fmt.Errorf("failed to create commit with readme: %w", err)
		}

		if _, err := svc.snap.Snapshot(
			codebaseID, ws.ID,
			snapshots.ActionViewSync, // TODO: Dedicated action for this?
			snapshotter.WithOnTemporaryView(),
			snapshotter.WithMarkAsLatestInWorkspace(),
			snapshotter.WithOnExistingCommit(commitID),
			snapshotter.WithOnRepo(repo), // Re-use repo context
		); err != nil {
			return fmt.Errorf("failed to create snapshot: %w", err)
		}

		return nil
	}

	if err := svc.executorProvider.New().FileReadGitWrite(cb).ExecTrunk(codebaseID, "createWelcomeMessage"); err != nil {
		return fmt.Errorf("failed to create welcome snapshot: %w", err)
	}

	return nil
}

func (s *WorkspaceService) RemovePatches(ctx context.Context, allower *unidiff.Allower, ws *workspaces.Workspace, hunkIDs ...string) error {
	removePatches := vcs_workspace.Remove(s.logger, hunkIDs...)

	if ws.ViewID != nil {
		if err := s.executorProvider.New().Write(removePatches).ExecView(ws.CodebaseID, *ws.ViewID, "removePatches"); err != nil {
			return fmt.Errorf("failed to remove patches: %w", err)
		}
		return nil
	}

	if ws.LatestSnapshotID != nil {
		snapshot, err := s.snap.GetByID(ctx, *ws.LatestSnapshotID)
		if err != nil {
			return fmt.Errorf("failed to get snapshot: %w", err)
		}
		if err := s.executorProvider.New().
			Write(vcs_view.CheckoutSnapshot(snapshot)).
			Write(func(repo vcs.RepoWriter) error {
				if err := removePatches(repo); err != nil {
					return fmt.Errorf("failed to remove patches: %w", err)
				}

				if _, err := s.snap.Snapshot(
					ws.CodebaseID,
					ws.ID,
					snapshots.ActionFileUndoPatch,
					snapshotter.WithOnView(*repo.ViewID()),
					snapshotter.WithMarkAsLatestInWorkspace(),
					snapshotter.WithOnRepo(repo),
				); err != nil {
					return fmt.Errorf("failed to snapshot: %w", err)
				}

				return nil
			}).ExecTemporaryView(ws.CodebaseID, "removePatches"); err != nil {
			return fmt.Errorf("failed to remove patches: %w", err)
		}

		return nil
	}

	return fmt.Errorf("failed to remove patches: no view or snapshot")
}

func (s *WorkspaceService) HasConflicts(ctx context.Context, ws *workspaces.Workspace) (bool, error) {
	if ws.LatestSnapshotID == nil {
		// can not check for conflicts, have no snapshot
		return false, nil
	}

	snapshotBranchName := fmt.Sprintf("snapshot-" + *ws.LatestSnapshotID)

	var hasConflicts bool
	checkConflicts := func(repo vcs.RepoGitWriter) error {
		idx, err := repo.MergeBranches(snapshotBranchName, "sturdytrunk")
		if err != nil {
			return fmt.Errorf("failed to merge branches: %w", err)
		}
		defer idx.Free()

		hasConflicts = idx.HasConflicts()
		return nil
	}

	checkConflictsOnView := func(repo vcs.RepoGitWriter) error {
		// If sturdytrunk doesn't exist (such as when an empty repository has been imported), it's not conflicting
		if _, err := repo.BranchCommitID("sturdytrunk"); err != nil {
			return nil
		}

		if err := repo.FetchBranch(snapshotBranchName, "sturdytrunk"); err != nil {
			return fmt.Errorf("failed to fetch branch: %w", err)
		}

		return checkConflicts(repo)
	}

	checkConflictsOnTrunk := func(repo vcs.RepoGitWriter) error {
		// If sturdytrunk doesn't exist (such as when an empty repository has been imported), it's not conflicting
		if _, err := repo.BranchCommitID("sturdytrunk"); err != nil {
			return nil
		}

		return checkConflicts(repo)
	}

	if ws.ViewID == nil {
		snapshot, err := s.snap.GetByID(ctx, *ws.LatestSnapshotID)
		if err != nil {
			return false, fmt.Errorf("failed to get snapshot: %w", err)
		}
		if err := s.executorProvider.New().
			Write(vcs_view.CheckoutSnapshot(snapshot)).
			GitWrite(checkConflictsOnTrunk).
			ExecTrunk(ws.CodebaseID, "workspaceCheckIfConflictsWithSnapshot"); err != nil {
			return false, fmt.Errorf("failed to check if conflicts: %w", err)
		}
		return hasConflicts, nil
	} else {
		if err := s.executorProvider.New().
			GitWrite(checkConflictsOnView).
			ExecView(ws.CodebaseID, *ws.ViewID, "workspaceCheckIfConflictsOnView"); err != nil {
			if errors.Is(err, executor.ErrIsRebasing) {
				return false, nil
			}
			return false, err
		}
		return hasConflicts, nil
	}
}

// ArchiveWithChange is the same as Archive, but also marks workspacw with the change ID.
func (s *WorkspaceService) ArchiveWithChange(ctx context.Context, ws *workspaces.Workspace, change *changes.Change) error {
	return s.archive(ctx, ws, &change.ID)
}

// Archive archives a workspace. If there is a view connected to the workspace, it will be reconnected to a new workspace.
func (s *WorkspaceService) Archive(ctx context.Context, ws *workspaces.Workspace) error {
	return s.archive(ctx, ws, nil)
}

func (s *WorkspaceService) archive(ctx context.Context, ws *workspaces.Workspace, changeID *changes.ID) error {
	if ws.ArchivedAt != nil {
		return nil // noop
	}

	t := time.Now()
	ws.ArchivedAt = &t
	ws.UnarchivedAt = nil
	updateFields := []db.UpdateOption{db.SetArchivedAt(&t), db.SetUnarchivedAt(nil)}

	if changeID != nil {
		ws.ChangeID = changeID
		updateFields = append(updateFields, db.SetChangeID(changeID))
	}

	if err := s.workspaceWriter.UpdateFields(ctx, ws.ID, updateFields...); err != nil {
		return fmt.Errorf("failed to archive workspace: %w", err)
	}

	s.analyticsService.Capture(ctx, "workspace archived", analytics.CodebaseID(ws.CodebaseID),
		analytics.Property("workspace_id", ws.ID),
	)

	if ws.ViewID == nil {
		return nil
	}

	// if there is a view connected, move it to a new workspace

	view, err := s.viewService.GetByID(ctx, *ws.ViewID)
	if err != nil {
		return fmt.Errorf("failed to get view: %w", err)
	}

	newWorkspace, err := s.Create(ctx, CreateWorkspaceRequest{
		CodebaseID: ws.CodebaseID,
		UserID:     ws.UserID,
	})
	if err != nil {
		return fmt.Errorf("failed to create new workspace: %w", err)
	}

	if err := s.viewService.OpenWorkspace(ctx, view, newWorkspace); err != nil {
		return fmt.Errorf("failed to open workspace on view: %w", err)
	}

	// Send events that the codebase has been updated, list of workspaces has changed
	if err := s.eventsSender.Codebase(ws.CodebaseID, events.CodebaseUpdated, ws.CodebaseID.String()); err != nil {
		s.logger.Error("failed to send codebase event", zap.Error(err))
		// do not fail
	}

	return nil
}

func (s *WorkspaceService) Unarchive(ctx context.Context, ws *workspaces.Workspace) error {
	t := time.Now()
	ws.ArchivedAt = nil
	ws.UnarchivedAt = &t
	ws.ViewID = nil
	ws.ChangeID = nil
	if err := s.workspaceWriter.UpdateFields(ctx, ws.ID,
		db.SetArchivedAt(nil),
		db.SetUnarchivedAt(&t),
		db.SetViewID(nil),
		db.SetChangeID(nil),
	); err != nil {
		return fmt.Errorf("failed to unarchive workspace: %w", err)
	}
	s.analyticsService.Capture(ctx, "workspace unarchived", analytics.CodebaseID(ws.CodebaseID),
		analytics.Property("workspace_id", ws.ID),
	)

	// Send events that the codebase has been updated, list of workspaces has changed
	if err := s.eventsSender.Codebase(ws.CodebaseID, events.CodebaseUpdated, ws.CodebaseID.String()); err != nil {
		s.logger.Error("failed to send codebase event", zap.Error(err))
		// do not fail
	}

	return nil
}

var ErrNotAvailable = errors.New("this features is not available in this version of Sturdy")

func (s *WorkspaceService) Push(ctx context.Context, user *users.User, ws *workspaces.Workspace) error {
	return ErrNotAvailable
}

func (s *WorkspaceService) LandOnSturdyAndPushTracked(ctx context.Context, ws *workspaces.Workspace) error {
	return ErrNotAvailable
}
