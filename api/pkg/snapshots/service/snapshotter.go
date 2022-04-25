package service

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"getsturdy.com/api/pkg/analytics"
	service_analytics "getsturdy.com/api/pkg/analytics/service"
	"getsturdy.com/api/pkg/codebases"
	"getsturdy.com/api/pkg/events"
	eventsv2 "getsturdy.com/api/pkg/events/v2"
	"getsturdy.com/api/pkg/snapshots"
	db_snapshots "getsturdy.com/api/pkg/snapshots/db"
	vcs_snapshots "getsturdy.com/api/pkg/snapshots/vcs"
	db_suggestions "getsturdy.com/api/pkg/suggestions/db"
	"getsturdy.com/api/pkg/unidiff"
	"getsturdy.com/api/pkg/unidiff/lfs"
	db_view "getsturdy.com/api/pkg/view/db"
	vcs_view "getsturdy.com/api/pkg/view/vcs"
	db_workspaces "getsturdy.com/api/pkg/workspaces/db"
	"getsturdy.com/api/vcs"
	"getsturdy.com/api/vcs/executor"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

type SnapshotOptions struct {
	patchIDsFilter          *[]string
	revertCommitHeadBase    *[2]*string
	onTemporaryView         bool
	onView                  *string
	onRepo                  vcs.RepoReaderGitWriter
	onExistingCommit        *string
	markAsLatestInWorkspace bool
	withNoThrottle          bool
}

type SnapshotOption func(*SnapshotOptions)

func WithPatchIDsFilter(patchIDs []string) SnapshotOption {
	return func(opts *SnapshotOptions) {
		if opts.patchIDsFilter == nil {
			opts.patchIDsFilter = new([]string)
		}
		*opts.patchIDsFilter = append(*opts.patchIDsFilter, patchIDs...)
	}
}

func WithRevertDiff(head string, base *string) SnapshotOption {
	return func(opts *SnapshotOptions) {
		opts.revertCommitHeadBase = &[2]*string{&head, base}
	}
}

func WithOnTemporaryView() SnapshotOption {
	return func(opts *SnapshotOptions) {
		opts.onTemporaryView = true
	}
}

func WithOnView(viewID string) SnapshotOption {
	return func(opts *SnapshotOptions) {
		opts.onView = &viewID
	}
}

func WithOnRepo(repo vcs.RepoReaderGitWriter) SnapshotOption {
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

func WithNoThrottle() SnapshotOption {
	return func(opts *SnapshotOptions) {
		opts.withNoThrottle = true
	}
}

type Service struct {
	snapshotsRepo   db_snapshots.Repository
	workspaceReader db_workspaces.WorkspaceReader
	workspaceWriter db_workspaces.WorkspaceWriter
	viewRepo        db_view.Repository
	suggestionsRepo db_suggestions.Repository

	eventsSender     events.EventSender
	eventsSenderV2   *eventsv2.Publisher
	executorProvider executor.Provider
	logger           *zap.Logger

	analyticsService *service_analytics.Service
}

func New(
	snapshotsRepo db_snapshots.Repository,
	workspaceReader db_workspaces.WorkspaceReader,
	workspaceWriter db_workspaces.WorkspaceWriter,
	viewRepo db_view.Repository,
	suggestionsRepo db_suggestions.Repository,

	eventSender events.EventSender,
	eventsSenderV2 *eventsv2.Publisher,
	executorProvider executor.Provider,
	logger *zap.Logger,

	analyticsService *service_analytics.Service,
) *Service {
	return &Service{
		snapshotsRepo:   snapshotsRepo,
		workspaceReader: workspaceReader,
		workspaceWriter: workspaceWriter,
		viewRepo:        viewRepo,
		suggestionsRepo: suggestionsRepo,

		eventsSender:     eventSender,
		eventsSenderV2:   eventsSenderV2,
		executorProvider: executorProvider,
		logger:           logger.Named("GitSnapshotter"),

		analyticsService: analyticsService,
	}
}

func getSnapshotOptions(opts ...SnapshotOption) *SnapshotOptions {
	options := &SnapshotOptions{}
	for _, applyOption := range opts {
		applyOption(options)
	}
	return options
}

func (s *Service) GetByID(_ context.Context, id string) (*snapshots.Snapshot, error) {
	return s.snapshotsRepo.Get(id)
}

var (
	ErrCantSnapshotRebasing    = errors.New("can't snapshot, rebasing in progress")
	ErrCantSnapshotWrongBranch = errors.New("can't snapshot, unexpected branch")
)

//nolint:cyclop
func (s *Service) Snapshot(codebaseID codebases.ID, workspaceID string, action snapshots.Action, opts ...SnapshotOption) (*snapshots.Snapshot, error) {
	options := getSnapshotOptions(opts...)

	if !options.onTemporaryView && options.onView == nil {
		return nil, errors.New("either onTemporaryView or onView must be set")
	}
	if options.onTemporaryView && options.onView != nil {
		return nil, errors.New("onTemporaryView and onView are mutually exclusive")
	}
	if options.onRepo != nil && (options.onView == nil && options.onExistingCommit == nil) {
		return nil, errors.New("when onRepo is set, onView or onExistingCommit must also be set")
	}
	if options.onExistingCommit != nil && options.onRepo == nil {
		return nil, errors.New("when onExistingCommit is set, onRepo must also be set")
	}

	logger := s.logger.With(
		zap.Stringer("codebase_id", codebaseID),
		zap.String("workspace_id", workspaceID),
		zap.Bool("option_on_temporary_view", options.onTemporaryView),
		zap.Stringp("option_on_view", options.onView),
		zap.Stringer("snapshot_action", action),
	)

	// if this view is the authoritative view of the workspace, mark this snapshot as the latest one
	ws, err := s.workspaceReader.Get(workspaceID)
	if err != nil {
		return nil, fmt.Errorf("failed to get workspace: %w", err)
	}

	var latest *snapshots.Snapshot

	if options.onView != nil {
		// Find previous snapshot
		var err error
		latest, err = s.snapshotsRepo.LatestInWorkspace(context.TODO(), workspaceID)
		if err != nil && !errors.Is(err, sql.ErrNoRows) {
			return nil, err
		}

		_, err = s.suggestionsRepo.GetByWorkspaceID(context.TODO(), workspaceID)
		if err != nil && !errors.Is(err, sql.ErrNoRows) {
			return nil, err
		}
		isSuggesting := err == nil

		// Throttle view sync snapshots, at most once per 10s
		// The throttling does not apply to workspaces that are suggesting
		if latest != nil &&
			!isSuggesting &&
			!options.withNoThrottle &&
			(ws.DiffsCount != nil && *ws.DiffsCount > 0) && // Always snapshot if we have no recorded diff count, or if the count is 0
			action == snapshots.ActionViewSync &&
			latest.Action == snapshots.ActionViewSync &&
			latest.CreatedAt.After(time.Now().Add(-10*time.Second)) {
			logger.Info("throttle ActionViewSync snapshot", zap.Duration("since_last_duration", time.Since(latest.CreatedAt)))
			return nil, nil
		}
	}

	snapshotID := uuid.New().String()

	var (
		snapshotCommitID string
		diffsCount       int32
	)

	countDiffs := func(repo vcs.RepoReader) error {
		gitDiffs, err := repo.CurrentDiffNoIndex()
		if err != nil {
			return fmt.Errorf("can't get git diffs: %w", err)
		}

		ignoreLFS, err := lfs.NewIgnoreLfsSmudgedFilter(repo)
		if err != nil {
			return fmt.Errorf("could not smudge lfs files: %w", err)
		}

		differ := unidiff.NewUnidiff(unidiff.NewGitPatchReader(gitDiffs), s.logger).
			WithExpandedHunks().
			WithFilterFunc(ignoreLFS)

		diffs, err := differ.Decorate()
		if err != nil {
			return fmt.Errorf("can't decorate git diffs: %w", err)
		}
		diffsCount = int32(len(diffs))
		return nil
	}

	var snapshotOptions []vcs_snapshots.SnapshotOption
	if options.patchIDsFilter != nil {
		snapshotOptions = append(snapshotOptions, vcs_snapshots.WithPatchIDsFilter(*options.patchIDsFilter))
	}
	if options.revertCommitHeadBase != nil {
		snapshotOptions = append(snapshotOptions, vcs_snapshots.WithRevert(*options.revertCommitHeadBase[0], options.revertCommitHeadBase[1]))
	}

	if options.onTemporaryView && options.onExistingCommit != nil && options.onRepo != nil {
		var err error
		snapshotCommitID, err = vcs_snapshots.SnapshotOnExistingCommit(options.onRepo, snapshotID, *options.onExistingCommit)
		if err != nil {
			return nil, err
		}
	} else if options.onRepo != nil && options.onView != nil {
		if err := countDiffs(options.onRepo); err != nil {
			return nil, err
		}

		var err error
		snapshotCommitID, err = vcs_snapshots.SnapshotOnViewRepo(s.logger, options.onRepo, codebaseID, snapshotID, snapshotOptions...)
		if err != nil {
			return nil, err
		}
	} else if options.onRepo == nil {
		// Run in a new executor
		exec := s.executorProvider.New()
		if !options.onTemporaryView {
			exec = exec.AssertBranchName(workspaceID)
		}
		var err error

		if options.revertCommitHeadBase != nil {
			// Reverting snapshot
			exec = exec.Write(func(repo vcs.RepoWriter) error {
				commitID, err := vcs_snapshots.SnapshotOnViewRepoWithRevert(repo, s.logger, snapshotID, snapshotOptions...)
				if err != nil {
					return err
				}
				snapshotCommitID = commitID
				return nil
			})

			// TODO: this is not true for reverts
			// snapshot on trunk is basically a copy of a commit => no diffs
			diffsCount = 0
		} else {
			// Normal snapshot
			exec = exec.Read(countDiffs)
			exec = exec.FileReadGitWrite(func(repo vcs.RepoReaderGitWriter) error {
				commitID, err := vcs_snapshots.SnapshotOnViewRepo(s.logger, repo, codebaseID, snapshotID, snapshotOptions...)
				if err != nil {
					return err
				}
				snapshotCommitID = commitID
				return nil
			})
		}

		if options.onTemporaryView {
			err = exec.ExecTemporaryView(codebaseID, "snapshotOnTemporaryView")
		} else {
			err = exec.ExecView(codebaseID, *options.onView, "snapshotOnView")
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
		ID:          snapshotID,
		CommitSHA:   snapshotCommitID,
		CreatedAt:   time.Now(),
		WorkspaceID: workspaceID,
		CodebaseID:  codebaseID,
		Action:      action,
		DiffsCount:  &diffsCount,
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
		isAuthoritativeView := ws.ViewID != nil && *ws.ViewID == *options.onView

		// If authoritative view, or explicitly asked to mark this as the latest snapshot
		if isAuthoritativeView || options.markAsLatestInWorkspace {
			ws.SetSnapshot(snap)
			if err := s.workspaceWriter.UpdateFields(context.TODO(), ws.ID,
				db_workspaces.SetDiffsCount(snap.DiffsCount),
				db_workspaces.SetLatestSnapshotID(&snap.ID),
			); err != nil {
				return nil, fmt.Errorf("failed to update workspace: %w", err)
			}

			// workspace updated
			if err := s.eventsSender.Workspace(ws.ID, events.WorkspaceUpdated, ws.ID); err != nil {
				s.logger.Error("failed to workspace updated event", zap.Error(err))
				// do not fail
			}
		}

		if isAuthoritativeView {
			if err := s.sendViewEvents(workspaceID, *options.onView); err != nil {
				return nil, err
			}
		}
	}

	s.analyticsService.CaptureUser(context.TODO(), ws.UserID, "created snapshot",
		analytics.CodebaseID(ws.CodebaseID),
		analytics.Property("diffs_count", diffsCount),
	)

	return snap, nil
}

func (s *Service) sendViewEvents(workspaceID, viewID string) error {
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

		ctx := context.Background()

		if err := s.eventsSenderV2.ViewUpdated(ctx, eventsv2.Codebase(ownerView.CodebaseID), ownerView); err != nil {
			s.logger.Error("failed to send view updated event", zap.Error(err))
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

func (s *Service) Diffs(ctx context.Context, snapshotID string, oo ...DiffsOption) ([]unidiff.FileDiff, error) {
	snapshot, err := s.snapshotsRepo.Get(snapshotID)
	if err != nil {
		return nil, fmt.Errorf("could not get snapshot: %w", err)
	}
	return s.diffs(ctx, snapshot, oo...)
}

func (s *Service) diffs(ctx context.Context, snapshot *snapshots.Snapshot, oo ...DiffsOption) ([]unidiff.FileDiff, error) {
	options := getDiffOptions(oo...)

	var diffs []unidiff.FileDiff
	if err := s.executorProvider.New().GitRead(func(repo vcs.RepoGitReader) error {
		snapParent, err := repo.GetCommitParents(snapshot.CommitSHA)
		if err != nil {
			return fmt.Errorf("failed to get commit parents: %w", err)
		}
		if len(snapParent) != 1 {
			return fmt.Errorf("unexpected number of snapshot parents: %d, expected %d", len(snapParent), 1)
		}

		gitDiffs, err := repo.DiffCommits(snapParent[0], snapshot.CommitSHA)
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
func (s *Service) Copy(ctx context.Context, snapshotID string, oo ...CopyOption) (*snapshots.Snapshot, error) {
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
		ViewID:      snapshot.ViewID,
		WorkspaceID: snapshot.WorkspaceID,
		Action:      snapshot.Action,
		DiffsCount:  snapshot.DiffsCount,
	}

	if err := s.executorProvider.New().
		Write(vcs_view.CheckoutBranch(snapshot.WorkspaceID)).
		Write(func(repo vcs.RepoWriter) error {
			if err := repo.ApplyPatchesToWorkdir(patches); err != nil {
				return fmt.Errorf("failed to apply patches to workdir: %w", err)
			}

			commitID, err := vcs_snapshots.SnapshotOnViewRepo(s.logger, repo, newSnapshot.CodebaseID, newSnapshot.ID)
			if err != nil {
				return fmt.Errorf("failed to snapshot on view repo: %w", err)
			}
			newSnapshot.CommitSHA = commitID
			return nil
		}).ExecTemporaryView(snapshot.CodebaseID, "copySnapshot"); err != nil {
		return nil, fmt.Errorf("failed to copy snapshot: %w", err)
	}

	if err := s.snapshotsRepo.Create(newSnapshot); err != nil {
		return nil, fmt.Errorf("failed to create new snapshot: %w", err)
	}

	return newSnapshot, nil
}

func (s *Service) Restore(snap *snapshots.Snapshot, viewRepo vcs.RepoWriter) error {
	if err := vcs_snapshots.RestoreRepo(s.logger, viewRepo, snap.ID, snap.CommitSHA); err != nil {
		return fmt.Errorf("failed to restore: %w", err)
	}
	return nil
}

func (s *Service) GetByCommitSHA(ctx context.Context, sha string) (*snapshots.Snapshot, error) {
	return s.snapshotsRepo.GetByCommitSHA(ctx, sha)
}
