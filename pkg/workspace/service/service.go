package service

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	"mash/pkg/change/decorate"
	"mash/pkg/change/message"
	service_change "mash/pkg/change/service"
	change_vcs "mash/pkg/change/vcs"
	workers_ci "mash/pkg/ci/workers"
	service_comments "mash/pkg/comments/service"
	service_github "mash/pkg/github/service"
	db_review "mash/pkg/review/db"
	"mash/pkg/snapshots"
	"mash/pkg/snapshots/snapshotter"
	worker_snapshots "mash/pkg/snapshots/worker"
	"mash/pkg/unidiff"
	"mash/pkg/unidiff/lfs"
	user_db "mash/pkg/user/db"
	"mash/pkg/view/events"
	vcs_view "mash/pkg/view/vcs"
	"mash/pkg/workspace"
	"mash/pkg/workspace/activity"
	"mash/pkg/workspace/activity/sender"
	"mash/pkg/workspace/db"
	vcs_workspace "mash/pkg/workspace/vcs"
	"mash/vcs"
	"mash/vcs/executor"
	"mash/vcs/provider"

	"github.com/google/uuid"
	git "github.com/libgit2/git2go/v33"
	"github.com/posthog/posthog-go"
	"go.uber.org/zap"
)

type CreateWorkspaceRequest struct {
	UserID           string
	CodebaseID       string
	Name             string
	RevertChangeID   string
	ChangeID         string
	DraftDescription string
}

type Service interface {
	Create(CreateWorkspaceRequest) (*workspace.Workspace, error)
	GetByID(context.Context, string) (*workspace.Workspace, error)
	LandChange(ctx context.Context, ws *workspace.Workspace, patchIDs []string, diffOptions ...vcs.DiffOption) error
	CreateWelcomeWorkspace(codebaseID, userID, codebaseName string) error
	Diffs(context.Context, string, ...DiffsOption) ([]unidiff.FileDiff, bool, error)
	CopyPatches(ctx context.Context, src, dist *workspace.Workspace, opts ...CopyPatchesOption) error
	RemovePatches(context.Context, *unidiff.Allower, *workspace.Workspace, ...string) error
	HasConflicts(context.Context, *workspace.Workspace) (bool, error)
}

type workspaceService struct {
	logger        *zap.Logger
	postHogClient posthog.Client

	workspaceWriter db.WorkspaceWriter
	workspaceReader db.WorkspaceReader

	userRepo   user_db.Repository
	reviewRepo db_review.ReviewRepository

	commentService *service_comments.Service
	changeService  *service_change.Service
	gitHubService  *service_github.Service

	activitySender   sender.ActivitySender
	eventsSender     events.EventSender
	snapshotterQueue worker_snapshots.Queue
	executorProvider executor.Provider
	snap             snapshotter.Snapshotter
	buildQueue       *workers_ci.BuildQueue
}

func New(
	logger *zap.Logger,
	postHogClient posthog.Client,

	workspaceWriter db.WorkspaceWriter,
	workspaceReader db.WorkspaceReader,

	userRepo user_db.Repository,
	reviewRepo db_review.ReviewRepository,

	commentsService *service_comments.Service,
	changeService *service_change.Service,
	gitHubService *service_github.Service,

	activitySender sender.ActivitySender,
	executorProvider executor.Provider,
	eventsSender events.EventSender,
	snapshotterQueue worker_snapshots.Queue,
	snap snapshotter.Snapshotter,
	buildQueue *workers_ci.BuildQueue,
) Service {
	return &workspaceService{
		logger:        logger,
		postHogClient: postHogClient,

		workspaceWriter: workspaceWriter,
		workspaceReader: workspaceReader,

		userRepo:   userRepo,
		reviewRepo: reviewRepo,

		commentService: commentsService,
		changeService:  changeService,
		gitHubService:  gitHubService,

		activitySender:   activitySender,
		executorProvider: executorProvider,
		eventsSender:     eventsSender,
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

func (s *workspaceService) Diffs(ctx context.Context, workspaceID string, oo ...DiffsOption) ([]unidiff.FileDiff, bool, error) {
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

func (s *workspaceService) diffsFromSnapshot(ctx context.Context, ws *workspace.Workspace, options *DiffsOptions) ([]unidiff.FileDiff, error) {
	if ws.LatestSnapshotID == nil {
		return nil, nil
	}

	snapshotOptions := []snapshotter.DiffsOption{}
	if options.Allower != nil {
		snapshotOptions = append(snapshotOptions, snapshotter.WithAllower(options.Allower))
	}

	return s.snap.Diffs(ctx, *ws.LatestSnapshotID, snapshotOptions...)
}

func (s *workspaceService) diffsFromView(ctx context.Context, ws *workspace.Workspace, options *DiffsOptions) ([]unidiff.FileDiff, bool, error) {
	var diffs []unidiff.FileDiff

	isRebasing := false
	if err := s.executorProvider.New().AssertBranchName(ws.ID).AllowRebasingState().Read(func(repo vcs.RepoReader) error {
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

func (s *workspaceService) GetByID(ctx context.Context, id string) (*workspace.Workspace, error) {
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

func (s *workspaceService) CopyPatches(ctx context.Context, dist, src *workspace.Workspace, opts ...CopyPatchesOption) error {
	if src.CodebaseID != dist.CodebaseID {
		return fmt.Errorf("source and destination codebases must be the same")
	}

	if dist.ViewID != nil {
		// TODO
		return fmt.Errorf("copying to active workspace is not supported")
	}

	options := getCopyPatchOptions(opts...)
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
		dist.LatestSnapshotID = &snapshot.ID
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
		dist.LatestSnapshotID = &snapshot.ID
	} else {
		// if we don't need to copy patches, re-use the existing snapshot
		dist.LatestSnapshotID = src.LatestSnapshotID
	}

	if err := s.workspaceWriter.Update(dist); err != nil {
		return fmt.Errorf("failed to update workspace: %w", err)
	}

	return nil
}

func (s *workspaceService) Create(req CreateWorkspaceRequest) (*workspace.Workspace, error) {
	t := time.Now()
	ws := workspace.Workspace{
		ID:               uuid.New().String(),
		UserID:           req.UserID,
		CodebaseID:       req.CodebaseID,
		CreatedAt:        &t,
		DraftDescription: req.DraftDescription,
	}

	if len(req.Name) > 0 {
		ws.Name = &req.Name
	} else {
		user, err := s.userRepo.Get(req.UserID)
		if err != nil {
			s.logger.Error("failed to get name of user when creating workspace", zap.Error(err))
		} else {
			name := fmt.Sprintf("%s's Workspace", user.Name)
			ws.Name = &name
		}
	}

	if err := s.executorProvider.New().Git(func(repo vcs.Repo) error {
		// Ensure codebase status
		if err := EnsureCodebaseStatus(repo); err != nil {
			return err
		}

		if req.RevertChangeID != "" {
			// Create workspace at the change that we want to revert
			if err := vcs_workspace.CreateAtChange(repo, ws.ID, req.RevertChangeID); err != nil {
				return fmt.Errorf("failed to create workspace at change: %w", err)
			}
		} else if req.ChangeID != "" {
			// Create workspace at any change
			if err := vcs_workspace.CreateAtChange(repo, ws.ID, req.ChangeID); err != nil {
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

	// Add the reverted changes to a snapshot
	if req.RevertChangeID != "" {
		if snapshot, err := s.snap.Snapshot(
			ws.CodebaseID,
			ws.ID,
			snapshots.ActionChangeReverted,
			snapshotter.WithOnTrunk(),
			snapshotter.WithRevertCommitID(req.RevertChangeID),
		); err != nil {
			return nil, fmt.Errorf("failed to create snapshot for revert: %w", err)
		} else {
			ws.LatestSnapshotID = &snapshot.ID
		}
	}

	if err := s.workspaceWriter.Create(ws); err != nil {
		return nil, fmt.Errorf("failed to write workspace to db: %w", err)
	}

	if err := s.postHogClient.Enqueue(posthog.Capture{
		DistinctId: req.UserID,
		Event:      "create workspace",
		Properties: posthog.NewProperties().
			Set("codebase_id", req.CodebaseID).
			Set("name", ws.Name).
			Set("at_existing_change", req.ChangeID != ""),
	}); err != nil {
		s.logger.Error("posthog failed", zap.Error(err))
	}

	return &ws, nil
}

func (s *workspaceService) LandChange(ctx context.Context, ws *workspace.Workspace, patchIDs []string, diffOpts ...vcs.DiffOption) error {
	gitHubRepository, err := s.gitHubService.GetRepositoryByCodebaseID(ctx, ws.CodebaseID)
	switch {
	case err == nil, errors.Is(err, sql.ErrNoRows):
	default:
		return fmt.Errorf("failed to get gitHubRepository: %w", err)
	}

	if gitHubRepository != nil && gitHubRepository.IntegrationEnabled && gitHubRepository.GitHubSourceOfTruth {
		return fmt.Errorf("landing disallowed when a github integration exists for codebase (github is source of truth)")
	}

	user, err := s.userRepo.Get(ws.UserID)
	if err != nil {
		return fmt.Errorf("failed to get user: %w", err)
	}

	cleanCommitMessage := message.CommitMessage(ws.DraftDescription)

	changeMeta := decorate.ChangeMetadata{
		Description: cleanCommitMessage,
		UserID:      user.ID,
	}
	if ws.ViewID != nil {
		changeMeta.ViewID = *ws.ViewID
	}
	gitCommitMessage := changeMeta.ToCommitMessage()

	signature := git.Signature{
		Name:  user.Name,
		Email: user.Email,
		When:  time.Now(),
	}

	var createdCommitID string

	var fromViewPushFunc func(vcs.Repo) error
	var fromSnapshotPushFunc func() error

	if ws.ViewID != nil {
		if err := s.executorProvider.New().Write(func(viewRepo vcs.RepoWriter) error {
			var err error
			createdCommitID, fromViewPushFunc, err = change_vcs.CreateAndLandFromView(
				viewRepo,
				s.logger,
				ws.CodebaseID,
				ws.ID,
				*ws.ViewID,
				patchIDs,
				gitCommitMessage,
				signature,
				diffOpts...,
			)
			if err != nil {
				return fmt.Errorf("failed to create and land from view: %w", err)
			}
			return nil
		}).ExecView(ws.CodebaseID, *ws.ViewID, "landChangeCreateAndLandFromView"); err != nil {
			return fmt.Errorf("failed to share from view: %w", err)
		}
	} else {
		if ws.LatestSnapshotID == nil {
			return fmt.Errorf("the workspace has no snapshot")
		}
		if err := s.executorProvider.New().Schedule(func(repoProvider provider.RepoProvider) error {
			var err error
			createdCommitID, fromSnapshotPushFunc, err = change_vcs.CreateAndLandFromSnapshot(
				repoProvider,
				s.logger,
				ws.CodebaseID,
				ws.ID,
				*ws.LatestSnapshotID,
				patchIDs,
				gitCommitMessage,
				signature,
				diffOpts...,
			)
			if err != nil {
				return fmt.Errorf("failed to create and land from snapshot: %w", err)
			}
			return nil
		}).ExecTrunk(ws.CodebaseID, "landChangeCreateAndLandFromSnapshot"); err != nil {
			return fmt.Errorf("failed to create and land from snaphsot: %w", err)
		}
	}

	cleanCommitMessageTitle := strings.Split(cleanCommitMessage, "\n")[0]
	change, err := s.changeService.Create(ctx, ws, createdCommitID, cleanCommitMessageTitle)
	if err != nil {
		return fmt.Errorf("failed to create change: %w", err)
	}

	if ws.ViewID != nil {
		executor := s.executorProvider.New().Git(func(viewRepo vcs.Repo) error {
			if err := fromViewPushFunc(viewRepo); err != nil {
				return fmt.Errorf("failed to push the landed result: %w", err)
			}
			return nil
		})
		if err := executor.ExecView(ws.CodebaseID, *ws.ViewID, "landChangePush"); err != nil {
			return fmt.Errorf("failed to push the landed result from view: %w", err)
		}
	} else {
		executor := s.executorProvider.New().Git(func(repo vcs.Repo) error {
			if err := fromSnapshotPushFunc(); err != nil {
				return fmt.Errorf("failed to push the landed result: %w", err)
			}
			return nil
		})
		if err := executor.ExecTrunk(ws.CodebaseID, "landChangePush"); err != nil {
			return fmt.Errorf("failed to push the landed result from trunk: %w", err)
		}
		// The snapshot if merged and can't be merged again.
		ws.LatestSnapshotID = nil
	}

	if err := s.postHogClient.Enqueue(posthog.Capture{
		DistinctId: ws.UserID,
		Event:      "created change",
		Properties: posthog.NewProperties().
			Set("codebase_id", ws.CodebaseID).
			Set("workspace_id", ws.ID).
			Set("commit_id", createdCommitID).
			Set("change_id", change.ID),
	}); err != nil {
		s.logger.Error("posthog failed", zap.Error(err))
	}

	if err := s.reviewRepo.DismissAllInWorkspace(ctx, ws.ID); err != nil {
		return fmt.Errorf("failed to dismiss all reviews: %w", err)
	}

	if ws.ViewID != nil {
		if err := s.snapshotterQueue.Enqueue(ctx, ws.CodebaseID, *ws.ViewID, ws.ID, []string{"."}, snapshots.ActionChangeLand); err != nil {
			return fmt.Errorf("failed to enqueue snapshot: %w", err)
		}

		if err := s.eventsSender.Codebase(ws.CodebaseID, events.ViewUpdated, *ws.ViewID); err != nil {
			return fmt.Errorf("failed to send view updated event: %w", err)
		}
	}

	if gitHubRepository != nil && gitHubRepository.IntegrationEnabled && !gitHubRepository.GitHubSourceOfTruth {
		if s.gitHubService.Push(ctx, gitHubRepository, change); err != nil {
			return fmt.Errorf("failed to push to github: %w", err)
		}
	}

	// Update workspace
	now := time.Now()
	ws.LastLandedAt = &now
	ws.UpdatedAt = &now
	ws.DraftDescription = ""
	ws.HeadCommitID = nil
	if err := s.workspaceWriter.Update(ws); err != nil {
		return fmt.Errorf("failed to update workspace: %w", err)
	}

	// Send event that the workspace has been updated
	if err := s.eventsSender.Workspace(ws.ID, events.WorkspaceUpdated, ws.ID); err != nil {
		s.logger.Error("failed to send workspace event", zap.Error(err))
	}

	// Clear 'up to date' cache for all workspaces
	if err := s.workspaceWriter.UnsetUpToDateWithTrunkForAllInCodebase(ws.CodebaseID); err != nil {
		return fmt.Errorf("failed to unset up_to_date_with_trunk: %w", err)
	}

	if err := s.postHogClient.Enqueue(posthog.Capture{
		DistinctId: ws.UserID,
		Event:      "landed changes",
		Properties: posthog.NewProperties().
			Set("codebase_id", ws.CodebaseID).
			Set("workspace_id", ws.ID).
			Set("change_id", change.ID),
	}); err != nil {
		s.logger.Error("posthog failed", zap.Error(err))
	}

	if err := s.commentService.MoveCommentsFromWorkspaceToChange(ctx, ws.ID, change.ID); err != nil {
		return fmt.Errorf("failed to move comments from workspace to change: %w", err)
	}

	// Create activity
	if err := s.activitySender.Codebase(ctx, ws.CodebaseID, ws.ID, ws.UserID, activity.WorkspaceActivityTypeCreatedChange, string(change.ID)); err != nil {
		return fmt.Errorf("failed to create workspace activity: %w", err)
	}

	// Send events that the codebase has been updated
	if err := s.eventsSender.Codebase(ws.CodebaseID, events.CodebaseUpdated, ws.CodebaseID); err != nil {
		s.logger.Error("failed to send codebase event", zap.Error(err))
	}

	if err := s.eventsSender.Workspace(ws.ID, events.WorkspaceUpdatedSnapshot, ws.ID); err != nil {
		s.logger.Error("failed to send workspace event", zap.Error(err))
	}

	if err := s.buildQueue.EnqueueChange(ctx, change); err != nil {
		s.logger.Error("failed to enqueue change", zap.Error(err))
	}

	return nil
}

func EnsureCodebaseStatus(repo vcs.Repo) error {
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
	<li><p><strong>Share</strong> this workspace to land the README on the trunk, and to make the file available to all collaborators</p></li>
</ul>

<p>Happy hacking!</p>
`

func (svc *workspaceService) CreateWelcomeWorkspace(codebaseID, userID, codebaseName string) error {
	readMeContents := strings.ReplaceAll(readMeTemplate, "__CODEBASE__NAME__", codebaseName)
	draftDescriptionContents := strings.ReplaceAll(draftDescriptionTemplate, "__CODEBASE__NAME__", codebaseName)

	ws, err := svc.Create(CreateWorkspaceRequest{
		CodebaseID:       codebaseID,
		UserID:           userID,
		Name:             "Add README",
		DraftDescription: draftDescriptionContents,
	})
	if err != nil {
		return fmt.Errorf("failed to create workspace: %w", err)
	}

	cb := func(repo vcs.RepoReader) error {
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
			snapshotter.WithOnTrunk(),
			snapshotter.WithMarkAsLatestInWorkspace(),
			snapshotter.WithOnExistingCommit(commitID),
			snapshotter.WithOnRepo(repo), // Re-use repo context
		); err != nil {
			return fmt.Errorf("failed to create snapshot: %w", err)
		}

		return nil
	}

	if err := svc.executorProvider.New().Read(cb).ExecTrunk(codebaseID, "createWelcomeMessage"); err != nil {
		return fmt.Errorf("failed to create welcome snapshot: %w", err)
	}

	return nil
}

func (s *workspaceService) RemovePatches(ctx context.Context, allower *unidiff.Allower, ws *workspace.Workspace, hunkIDs ...string) error {
	removePatches := vcs_workspace.Remove(s.logger, hunkIDs...)

	if ws.ViewID != nil {
		if err := s.executorProvider.New().Write(removePatches).ExecView(ws.CodebaseID, *ws.ViewID, "removePatches"); err != nil {
			return fmt.Errorf("failed to remove patches: %w", err)
		}
		return nil
	}

	if ws.LatestSnapshotID != nil {
		if err := s.executorProvider.New().Schedule(func(repoProvider provider.RepoProvider) error {
			repo, cancel, err := vcs_view.TemporaryViewFromSnapshot(repoProvider, ws.CodebaseID, ws.ID, *ws.LatestSnapshotID)
			if err != nil {
				return fmt.Errorf("failed to create temporary view: %w", err)
			}
			defer func() {
				if err := cancel(); err != nil {
					s.logger.Error("failed to cleanup temporary view", zap.Error(err))
				}
			}()
			defer repo.Free()

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
		}).ExecTrunk(ws.CodebaseID, "removePatches"); err != nil {
			return fmt.Errorf("failed to remove patches: %w", err)
		}

		return nil
	}

	return fmt.Errorf("failed to remove patches: no view or snapshot")
}

func (s *workspaceService) HasConflicts(ctx context.Context, ws *workspace.Workspace) (bool, error) {
	if ws.LatestSnapshotID == nil {
		// can not check for conflicts, have no snapshot
		return false, nil
	}

	var hasConflicts bool
	hasConflictsFunc := func(repo vcs.Repo) error {
		// If sturdytrunk doesn't exist (such as when an empty repository has been imported), it's not conflicting
		if _, err := repo.BranchCommitID("sturdytrunk"); err != nil {
			return nil
		}

		snapshotBranchName := fmt.Sprintf("snapshot-" + *ws.LatestSnapshotID)
		if err := repo.FetchBranch(snapshotBranchName, "sturdytrunk"); err != nil {
			return fmt.Errorf("failed to fetch branch: %w", err)
		}

		idx, err := repo.MergeBranches(snapshotBranchName, "sturdytrunk")
		if err != nil {
			return fmt.Errorf("failed to merge branches: %w", err)
		}
		defer idx.Free()

		hasConflicts = idx.HasConflicts()
		return nil
	}

	var err error
	if ws.ViewID == nil {
		err = s.executorProvider.New().Schedule(func(repoProvider provider.RepoProvider) error {
			repo, cancel, err := vcs_view.TemporaryViewFromSnapshot(repoProvider, ws.CodebaseID, ws.ID, *ws.LatestSnapshotID)
			if err != nil {
				return fmt.Errorf("failed to create temporary view: %w", err)
			}
			defer func() {
				if err := cancel(); err != nil {
					s.logger.Error("failed to cleanup temporary view", zap.Error(err))
				}
			}()
			defer repo.Free()
			return hasConflictsFunc(repo)
		}).ExecTrunk(ws.CodebaseID, "workspaceCheckIfConflicts")
	} else {
		err = s.executorProvider.New().Git(hasConflictsFunc).ExecView(ws.CodebaseID, *ws.ViewID, "workspaceCheckIfConflicts")
	}

	if err != nil {
		if errors.Is(err, executor.ErrIsRebasing) {
			return false, nil
		}
		return false, err
	}
	return hasConflicts, nil
}
