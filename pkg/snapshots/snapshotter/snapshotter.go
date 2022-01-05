package snapshotter

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"mash/pkg/snapshots"
	db_snapshots "mash/pkg/snapshots/db"
	vcs_snapshots "mash/pkg/snapshots/vcs"
	"mash/pkg/unidiff"
	db_view "mash/pkg/view/db"
	"mash/pkg/view/events"
	vcs_view "mash/pkg/view/vcs"
	db_workspace "mash/pkg/workspace/db"
	"mash/vcs"
	"mash/vcs/executor"
	"mash/vcs/provider"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

// todo: rename to snapshot service
type Snapshotter interface {
	Snapshot(codebaseID, workspaceID string, action snapshots.Action, options ...SnapshotOption) (*snapshots.Snapshot, error)
	Copy(ctx context.Context, snapshotID string, oo ...CopyOption) (*snapshots.Snapshot, error)
	Diffs(ctx context.Context, snapshotID string, oo ...DiffsOption) ([]unidiff.FileDiff, error)
	GetByID(context.Context, string) (*snapshots.Snapshot, error)
}

type SnapshotOptions struct {
	paths                   []string
	patchIDsFilter          *[]string
	revertCommitID          *string
	onTrunk                 bool
	onView                  *string
	onRepo                  vcs.RepoReader
	onExistingCommit        *string
	markAsLatestInWorkspace bool
}

type SnapshotOption func(*SnapshotOptions)

func WithPaths(paths []string) SnapshotOption {
	return func(opts *SnapshotOptions) {
		opts.paths = append(opts.paths, paths...)
	}
}

func WithPatchIDsFilter(patchIDs []string) SnapshotOption {
	return func(opts *SnapshotOptions) {
		if opts.patchIDsFilter == nil {
			opts.patchIDsFilter = new([]string)
		}
		*opts.patchIDsFilter = append(*opts.patchIDsFilter, patchIDs...)
	}
}

func WithRevertCommitID(commitID string) SnapshotOption {
	return func(opts *SnapshotOptions) {
		opts.revertCommitID = &commitID
	}
}

func WithOnTrunk() SnapshotOption {
	return func(opts *SnapshotOptions) {
		opts.onTrunk = true
	}
}

func WithOnView(viewID string) SnapshotOption {
	return func(opts *SnapshotOptions) {
		opts.onView = &viewID
	}
}

func WithOnRepo(repo vcs.RepoReader) SnapshotOption {
	return func(opts *SnapshotOptions) {
		opts.onRepo = repo
	}
}

func WithOnExistingCommit(commit string) SnapshotOption {
	return func(opts *SnapshotOptions) {
		opts.onExistingCommit = &commit
	}
}

func WithMarkAsLatestInWorkspace() SnapshotOption {
	return func(opts *SnapshotOptions) {
		opts.markAsLatestInWorkspace = true
	}
}

type snap struct {
	snapshotsRepo    db_snapshots.Repository
	workspaceReader  db_workspace.WorkspaceReader
	workspaceWriter  db_workspace.WorkspaceWriter
	viewRepo         db_view.Repository
	eventsSender     events.EventSender
	executorProvider executor.Provider
	logger           *zap.Logger
}

func NewGitSnapshotter(
	snapshotsRepo db_snapshots.Repository,
	workspaceReader db_workspace.WorkspaceReader,
	workspaceWriter db_workspace.WorkspaceWriter,
	viewRepo db_view.Repository,
	eventSender events.EventSender,
	executorProvider executor.Provider,
	logger *zap.Logger,
) Snapshotter {
	return &snap{
		snapshotsRepo:    snapshotsRepo,
		workspaceReader:  workspaceReader,
		workspaceWriter:  workspaceWriter,
		viewRepo:         viewRepo,
		eventsSender:     eventSender,
		executorProvider: executorProvider,
		logger:           logger.Named("GitSnapshotter"),
	}
}

func getSnapshotOptions(opts ...SnapshotOption) *SnapshotOptions {
	options := &SnapshotOptions{}
	for _, applyOption := range opts {
		applyOption(options)
	}
	return options
}

func (s *snap) GetByID(_ context.Context, id string) (*snapshots.Snapshot, error) {
	return s.snapshotsRepo.Get(id)
}

var (
	ErrCantSnapshotRebasing    = errors.New("can't snapshot, rebasing in progress")
	ErrCantSnapshotWrongBranch = errors.New("can't snapshot, unexpected branch")
)

func (s *snap) Snapshot(codebaseID, workspaceID string, action snapshots.Action, opts ...SnapshotOption) (*snapshots.Snapshot, error) {
	options := getSnapshotOptions(opts...)

	if !options.onTrunk && options.onView == nil {
		return nil, errors.New("either onTrunk or onView must be set")
	}
	if options.onTrunk && options.onView != nil {
		return nil, errors.New("onTrunk and onView are mutually exclusive")
	}
	if options.onRepo != nil && (options.onView == nil && options.onExistingCommit == nil) {
		return nil, errors.New("when onRepo is set, onView or onExistingCommit must also be set")
	}
	if options.onExistingCommit != nil && options.onRepo == nil {
		return nil, errors.New("when onExistingCommit is set, onRepo must also be set")
	}

	logger := s.logger.With(
		zap.String("codebase_id", codebaseID),
		zap.String("workspace_id", workspaceID),
		zap.Bool("option_on_trunk", options.onTrunk),
		zap.Stringp("option_on_view", options.onView),
		zap.Stringer("snapshot_action", action),
	)

	var latest *snapshots.Snapshot

	if options.onView != nil {
		// Find previous snapshot
		var err error
		latest, err = s.snapshotsRepo.LatestInView(*options.onView)
		if err != nil && !errors.Is(err, sql.ErrNoRows) {
			return nil, err
		}

		// Throttle view sync snapshots, at most once per 10s
		// The throttling does not apply to views that are suggesting
		if latest != nil &&
			action == snapshots.ActionViewSync &&
			latest.Action == snapshots.ActionViewSync &&
			latest.CreatedAt.After(time.Now().Add(-10*time.Second)) {
			logger.Info("throttle ActionViewSync snapshot", zap.Duration("since_last_duration", time.Since(latest.CreatedAt)))
			return nil, nil
		}
	}

	snapshotID := uuid.New().String()
	var snapshotCommitID string

	var snapshotOptions []vcs_snapshots.SnapshotOption
	if options.patchIDsFilter != nil {
		snapshotOptions = append(snapshotOptions, vcs_snapshots.WithPatchIDsFilter(*options.patchIDsFilter))
	}
	if options.revertCommitID != nil {
		snapshotOptions = append(snapshotOptions, vcs_snapshots.WithRevertCommitID(*options.revertCommitID))
	}

	if options.onTrunk && options.onExistingCommit != nil && options.onRepo != nil {
		var err error
		snapshotCommitID, err = vcs_snapshots.SnapshotOnExistingCommit(options.onRepo, snapshotID, *options.onExistingCommit)
		if err != nil {
			return nil, err
		}
	} else if options.onRepo != nil && options.onView != nil {
		var err error
		snapshotCommitID, err = vcs_snapshots.SnapshotOnViewRepo(s.logger, options.onRepo, codebaseID, snapshotID, snapshotOptions...)
		if err != nil {
			return nil, err
		}
	} else if options.onRepo == nil {
		// Run in a new executor
		exec := s.executorProvider.New()
		if !options.onTrunk {
			exec = exec.AssertBranchName(workspaceID)
		}
		var err error
		if options.onTrunk {
			err = exec.Git(func(repo vcs.Repo) error {
				commitID, err := vcs_snapshots.SnapshotOnTrunk(repo, workspaceID, snapshotID, snapshotOptions...)
				if err != nil {
					return err
				}
				snapshotCommitID = commitID
				return nil
			}).ExecTrunk(codebaseID, "snapshot")
		} else {
			err = exec.Read(func(repo vcs.RepoReader) error {
				commitID, err := vcs_snapshots.SnapshotOnViewRepo(s.logger, repo, codebaseID, snapshotID, snapshotOptions...)
				if err != nil {
					return err
				}
				snapshotCommitID = commitID
				return nil
			}).ExecView(codebaseID, *options.onView, "snapshot")
		}
		if errors.Is(err, executor.ErrUnexpectedBranch) {
			return nil, fmt.Errorf("%w: view is on unexpected branch (%s)", ErrCantSnapshotWrongBranch, err)
		}
		if errors.Is(err, executor.ErrIsRebasing) {
			return nil, fmt.Errorf("%w: view is rebasing", ErrCantSnapshotRebasing)
		}
		if err != nil {
			return nil, err
		}
	} else {
		return nil, fmt.Errorf("could not create snapshot, unrecognized combinations of options: %+v", options)
	}

	snap := &snapshots.Snapshot{
		ID:           snapshotID,
		CommitID:     snapshotCommitID,
		CreatedAt:    time.Now(),
		WorkspaceID:  &workspaceID,
		CodebaseID:   codebaseID,
		ChangedFiles: options.paths,
		Action:       action,
	}

	if options.onView != nil {
		snap.ViewID = *options.onView
	}

	if latest != nil {
		snap.PreviousSnapshotID = &latest.ID
	}

	if err := s.snapshotsRepo.Create(snap); err != nil {
		return nil, err
	}

	if options.onView != nil || options.markAsLatestInWorkspace {
		// if this view is the authoritative view of the workspace, mark this snapshot as the latest one
		ws, err := s.workspaceReader.Get(workspaceID)
		if err != nil {
			return nil, fmt.Errorf("failed to get workspace: %w", err)
		}

		isAuthoritativeView := ws.ViewID != nil && *ws.ViewID == *options.onView

		// If authoritative view, or explicitly asked to mark this as the latest snapshot
		if isAuthoritativeView || options.markAsLatestInWorkspace {
			ws.LatestSnapshotID = &snap.ID
			if err := s.workspaceWriter.Update(ws); err != nil {
				return nil, fmt.Errorf("failed to update workspace: %w", err)
			}
		}

		if isAuthoritativeView {
			if err := s.sendEvents(workspaceID, *options.onView); err != nil {
				return nil, err
			}
		}
	}

	return snap, nil
}

func (s *snap) sendEvents(workspaceID, viewID string) error {
	// If this is a _suggestion_, send events to the view it's making suggestions to
	ws, err := s.workspaceReader.Get(workspaceID)
	if errors.Is(err, sql.ErrNoRows) {
		return nil
	}
	if err != nil {
		return fmt.Errorf("could not get workspace: %w", err)
	}

	view, err := s.viewRepo.Get(viewID)
	if errors.Is(err, sql.ErrNoRows) {
		return nil
	}
	if err != nil {
		return fmt.Errorf("could not get view: %w", err)
	}

	if ws.UserID == view.UserID {
		return nil
	}
	// find the owners views
	ownerViews, err := s.viewRepo.ListByCodebaseAndUser(ws.CodebaseID, ws.UserID)
	if errors.Is(err, sql.ErrNoRows) {
		return nil
	}
	if err != nil {
		return fmt.Errorf("could not get workspace owner views: %w", err)
	}

	for _, ownerView := range ownerViews {
		if ownerView.WorkspaceID != workspaceID {
			continue
		}
		// TODO: this function must be called only once, it makes a db call
		if err := s.eventsSender.Codebase(ownerView.CodebaseID, events.ViewUpdated, ownerView.ID); err != nil {
			s.logger.Error("failed to send codebase event", zap.Error(err))
			// do not fail
		}
	}

	return nil
}

type DiffsOptions struct {
	Allower  *unidiff.Allower
	PatchIDs *[]string
}

type DiffsOption func(*DiffsOptions)

func DiffWithPatchIDs(patchIDs []string) DiffsOption {
	return func(options *DiffsOptions) {
		if options.PatchIDs == nil {
			options.PatchIDs = &patchIDs
		} else {
			*options.PatchIDs = append(*options.PatchIDs, patchIDs...)
		}
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

func (s *snap) Diffs(ctx context.Context, snapshotID string, oo ...DiffsOption) ([]unidiff.FileDiff, error) {
	snapshot, err := s.snapshotsRepo.Get(snapshotID)
	if err != nil {
		return nil, fmt.Errorf("could not get snapshot: %w", err)
	}
	return s.diffs(ctx, snapshot, oo...)
}

func (s *snap) diffs(ctx context.Context, snapshot *snapshots.Snapshot, oo ...DiffsOption) ([]unidiff.FileDiff, error) {
	options := getDiffOptions(oo...)

	var diffs []unidiff.FileDiff
	if err := s.executorProvider.New().Git(func(repo vcs.Repo) error {
		snapParent, err := repo.GetCommitParents(snapshot.CommitID)
		if err != nil {
			return fmt.Errorf("failed to get commit parents: %w", err)
		}
		if len(snapParent) != 1 {
			return fmt.Errorf("unexpected number of snapshot parents: %d, expected %d", len(snapParent), 1)
		}

		gitDiffs, err := repo.DiffCommits(snapParent[0], snapshot.CommitID)
		if err != nil {
			return fmt.Errorf("failed to get git diffs: %w", err)
		}
		defer gitDiffs.Free()

		differ := unidiff.NewUnidiff(unidiff.NewGitPatchReader(gitDiffs), s.logger).
			WithExpandedHunks()

		if options.Allower != nil {
			differ = differ.WithAllower(options.Allower)
		}

		if options.PatchIDs != nil {
			differ = differ.WithHunksFilter(*options.PatchIDs...)
		}

		hunkifiedDiff, err := differ.Decorate()
		if err != nil {
			return fmt.Errorf("failed to decorate diffs: %w", err)
		}
		diffs = hunkifiedDiff
		return nil
	}).ExecTrunk(snapshot.CodebaseID, "snapshotDiffs"); err != nil {
		return nil, fmt.Errorf("failed to get diffs from snapshot: %w", err)
	}
	return diffs, nil
}

type CopyOptions struct {
	PatchIDs *[]string
}

type CopyOption func(*CopyOptions)

func CopyWithPatchIDs(patchIDs []string) CopyOption {
	return func(options *CopyOptions) {
		if options.PatchIDs == nil {
			options.PatchIDs = &patchIDs
		} else {
			*options.PatchIDs = append(*options.PatchIDs, patchIDs...)
		}
	}
}

func getCopyOptions(oo ...CopyOption) *CopyOptions {
	options := &CopyOptions{}
	for _, opt := range oo {
		opt(options)
	}
	return options
}

// Copy creates a new snapshot from the given snapshot.
func (s *snap) Copy(ctx context.Context, snapshotID string, oo ...CopyOption) (*snapshots.Snapshot, error) {
	snapshot, err := s.snapshotsRepo.Get(snapshotID)
	if err != nil {
		return nil, fmt.Errorf("failed to get snapshot: %w", err)
	}

	diffOptions := []DiffsOption{}
	options := getCopyOptions(oo...)
	if options.PatchIDs != nil {
		diffOptions = append(diffOptions, DiffWithPatchIDs(*options.PatchIDs))
	}

	diffs, err := s.diffs(ctx, snapshot, diffOptions...)
	if err != nil {
		return nil, fmt.Errorf("failed to get diffs: %w", err)
	}

	patches := [][]byte{}
	for _, fd := range diffs {
		for _, hunk := range fd.Hunks {
			patches = append(patches, []byte(hunk.Patch))
		}
	}

	newSnapshot := &snapshots.Snapshot{
		ID:          uuid.New().String(),
		CreatedAt:   time.Now(),
		CodebaseID:  snapshot.CodebaseID,
		ViewID:      snapshot.CodebaseID,
		WorkspaceID: snapshot.WorkspaceID,
		Action:      snapshot.Action,
	}

	if err := s.executorProvider.New().Schedule(func(repoProvider provider.RepoProvider) error {
		repo, cancel, err := vcs_view.TemporaryView(repoProvider, newSnapshot.CodebaseID, *newSnapshot.WorkspaceID)
		if err != nil {
			return fmt.Errorf("failed to create temporary view: %w", err)
		}
		defer cancel()

		if err := repo.ApplyPatchesToWorkdir(patches); err != nil {
			return fmt.Errorf("failed to apply patches to workdir: %w", err)
		}

		commitID, err := vcs_snapshots.SnapshotOnViewRepo(s.logger, repo, newSnapshot.CodebaseID, newSnapshot.ID)
		if err != nil {
			return fmt.Errorf("failed to snapshot on view repo: %w", err)
		}
		newSnapshot.CommitID = commitID
		return nil
	}).ExecTrunk(snapshot.CodebaseID, "copySnapshot"); err != nil {
		return nil, fmt.Errorf("failed to copy snapshot: %w", err)
	}

	if err := s.snapshotsRepo.Create(newSnapshot); err != nil {
		return nil, fmt.Errorf("failed to create new snapshot: %w", err)
	}

	return newSnapshot, nil
}
