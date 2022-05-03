package service

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	git "github.com/libgit2/git2go/v33"

	"getsturdy.com/api/pkg/analytics"
	service_analytics "getsturdy.com/api/pkg/analytics/service"
	"getsturdy.com/api/pkg/codebases"
	"getsturdy.com/api/pkg/events"
	eventsv2 "getsturdy.com/api/pkg/events/v2"
	"getsturdy.com/api/pkg/snapshots"
	db_snapshots "getsturdy.com/api/pkg/snapshots/db"
	vcs_snapshots "getsturdy.com/api/pkg/snapshots/vcs"
	service_statuses "getsturdy.com/api/pkg/statuses/service"
	db_suggestions "getsturdy.com/api/pkg/suggestions/db"
	"getsturdy.com/api/pkg/unidiff"
	"getsturdy.com/api/pkg/unidiff/lfs"
	"getsturdy.com/api/pkg/users"
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
	withUser                *users.User
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

func WithUser(u *users.User) SnapshotOption {
	return func(opts *SnapshotOptions) {
		opts.withUser = u
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
	statusesService  *service_statuses.Service
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
	statusesService *service_statuses.Service,
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
		statusesService:  statusesService,
	}
}

func getSnapshotOptions(opts ...SnapshotOption) *SnapshotOptions {
	options := &SnapshotOptions{}
	for _, applyOption := range opts {
		applyOption(options)
	}
	return options
}

func (s *Service) GetByID(_ context.Context, id snapshots.ID) (*snapshots.Snapshot, error) {
	return s.snapshotsRepo.Get(id)
}

// Previous returns a snapshot that was made before this one.
func (s *Service) Previous(ctx context.Context, snapshot *snapshots.Snapshot) (*snapshots.Snapshot, error) {
	if snapshot.PreviousSnapshotID == nil {
		return nil, sql.ErrNoRows // todo: custom error?
	}
	return s.snapshotsRepo.Get(*snapshot.PreviousSnapshotID)
}

// Next returns a snapshot that references the given snapshot as a previous.
func (s *Service) Next(ctx context.Context, snapshot *snapshots.Snapshot) (*snapshots.Snapshot, error) {
	return s.snapshotsRepo.GetByPreviousSnapshotID(ctx, snapshot.ID)
}

var (
	ErrCantSnapshotRebasing    = errors.New("can't snapshot, rebasing in progress")
	ErrCantSnapshotWrongBranch = errors.New("can't snapshot, unexpected branch")
)

func (s *Service) Delete(ctx context.Context, snapshot *snapshots.Snapshot) error {
	now := time.Now()
	snapshot.DeletedAt = &now
	if err := s.snapshotsRepo.Update(snapshot); err != nil {
		return fmt.Errorf("can't update snapshot: %w", err)
	}
	return nil
}

//nolint:cyclop
func (s *Service) Snapshot(ctx context.Context, codebaseID codebases.ID, workspaceID string, action snapshots.Action, opts ...SnapshotOption) (*snapshots.Snapshot, error) {
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

	t0 := time.Now()
	snapshotID := snapshots.ID(uuid.New().String())

	logger := s.logger.With(
		zap.Stringer("snapshot_id", snapshotID),
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
	if ws.LatestSnapshotID != nil {
		previousSnapshot, err := s.snapshotsRepo.Get(*ws.LatestSnapshotID)
		if err != nil && !errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("failed to get latest snapshot: %w", err)
		}
		latest = previousSnapshot
	} else {
		previousSnapshot, err := s.snapshotsRepo.LatestInWorkspace(ctx, workspaceID)
		if err != nil && !errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("failed to get latest snapshot: %w", err)
		}
		latest = previousSnapshot
	}

	if options.onView != nil {
		if _, err := s.suggestionsRepo.GetByWorkspaceID(ctx, workspaceID); err != nil && !errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("failed to get suggestions: %w", err)
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

	var (
		snapshotCommitSHA string
		diffsCount        int32
		sameAsBefore      bool
	)

	compareTreeIDs := func(latest *snapshots.Snapshot, newSHA string) func(vcs.RepoGitReader) error {
		if latest == nil {
			return func(r vcs.RepoGitReader) error { return nil }
		}
		return func(repo vcs.RepoGitReader) error {
			newCommit, err := repo.Commit(newSHA)
			if err != nil {
				return fmt.Errorf("can't get new commit %s: %w", snapshotCommitSHA, err)
			}

			previousCommit, err := repo.Commit(latest.CommitSHA)
			if errors.Is(err, vcs.ErrNotFound) {
				return nil
			} else if err != nil {
				return fmt.Errorf("can't get previous commit %s: %w", latest.CommitSHA, err)
			}

			nt, err := newCommit.Tree()
			if err != nil {
				return fmt.Errorf("can't get new commit tree: %w", err)
			}

			pt, err := previousCommit.Tree()
			if err != nil {
				return fmt.Errorf("can't get previous commit tree: %w", err)
			}

			sameAsBefore = nt.Id().String() == pt.Id().String()

			return nil
		}
	}

	countDiffs := func(repo vcs.RepoReader) error {

		t0 := time.Now()

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

		logger.Info("counted diffs", zap.Duration("duration", time.Since(t0)))

		return nil
	}

	gitSignature := git.Signature{
		Name:  "Sturdy",
		Email: "support@getsturdy.com",
		When:  time.Now(),
	}
	if options.withUser != nil {
		gitSignature.Name = options.withUser.Name
		gitSignature.Email = options.withUser.Email
	}

	var snapshotOptions []vcs_snapshots.SnapshotOption
	if options.patchIDsFilter != nil {
		snapshotOptions = append(snapshotOptions, vcs_snapshots.WithPatchIDsFilter(*options.patchIDsFilter))
	}
	if options.revertCommitHeadBase != nil {
		snapshotOptions = append(snapshotOptions, vcs_snapshots.WithRevert(*options.revertCommitHeadBase[0], options.revertCommitHeadBase[1]))
	}

	// TODO: add the workspace name to the commit message
	snapshotOptions = append(snapshotOptions, vcs_snapshots.WithCommitMessage("Snapshot of "+workspaceID))

	if options.onTemporaryView && options.onExistingCommit != nil && options.onRepo != nil {
		var err error
		snapshotCommitSHA, err = vcs_snapshots.SnapshotOnExistingCommit(options.onRepo, snapshotID, *options.onExistingCommit)
		if err != nil {
			return nil, fmt.Errorf("can't snapshot on existing commit: %w", err)
		}

		if err := compareTreeIDs(latest, snapshotCommitSHA)(options.onRepo); err != nil {
			return nil, fmt.Errorf("can't compare trees: %w", err)
		}
	} else if options.onRepo != nil && options.onView != nil {
		if err := countDiffs(options.onRepo); err != nil {
			return nil, fmt.Errorf("can't count diffs: %w", err)
		}

		var err error
		snapshotCommitSHA, err = vcs_snapshots.SnapshotOnViewRepo(ctx, s.logger, options.onRepo, codebaseID, snapshotID, gitSignature, snapshotOptions...)
		if err != nil {
			return nil, fmt.Errorf("failed to snapshot on view repo: %w", err)
		}
		if err := compareTreeIDs(latest, snapshotCommitSHA)(options.onRepo); err != nil {
			return nil, fmt.Errorf("can't compare trees: %w", err)
		}
	} else if options.onRepo == nil {
		// Run in a new executor
		exec := s.executorProvider.New()
		if options.revertCommitHeadBase != nil {
			// Reverting snapshot
			exec = exec.Write(func(repo vcs.RepoWriter) error {
				commitID, err := vcs_snapshots.SnapshotOnViewRepoWithRevert(repo, s.logger, snapshotID, snapshotOptions...)
				if err != nil {
					return fmt.Errorf("failed to snapshot on view repo: %w", err)
				}
				snapshotCommitSHA = commitID
				if err := compareTreeIDs(latest, commitID)(repo); err != nil {
					return fmt.Errorf("can't compare trees: %w", err)
				}
				return nil
			})

			// TODO: this is not true for reverts
			// snapshot on trunk is basically a copy of a commit => no diffs
			diffsCount = 0
		} else {
			// Normal snapshot
			exec = exec.Read(countDiffs)
			exec = exec.FileReadGitWrite(func(repo vcs.RepoReaderGitWriter) error {
				commitID, err := vcs_snapshots.SnapshotOnViewRepo(ctx, s.logger, repo, codebaseID, snapshotID, gitSignature, snapshotOptions...)
				if err != nil {
					return fmt.Errorf("failed to snapshot on view repo: %w", err)
				}
				snapshotCommitSHA = commitID

				if err := compareTreeIDs(latest, commitID)(repo); err != nil {
					return fmt.Errorf("can't compare trees: %w", err)
				}
				return nil
			})
		}

		var err error
		if options.onTemporaryView {
			err = exec.Write(vcs_view.CheckoutBranch(workspaceID)).ExecTemporaryView(codebaseID, "snapshotOnTemporaryView")
		} else {
			err = exec.ExecView(codebaseID, *options.onView, "snapshotOnView")
		}

		if errors.Is(err, executor.ErrUnexpectedBranch) {
			return nil, fmt.Errorf("%w: view is on unexpected branch (%s)", ErrCantSnapshotWrongBranch, err)
		} else if errors.Is(err, executor.ErrIsRebasing) {
			return nil, fmt.Errorf("%w: view is rebasing", ErrCantSnapshotRebasing)
		} else if err != nil {
			return nil, fmt.Errorf("can't snapshot: %w", err)
		}
	} else {
		return nil, fmt.Errorf("could not create snapshot, unrecognized combinations of options: %+v", options)
	}

	isLatestDeleted := latest != nil && latest.IsDeleted()
	if shouldSnapshot := !sameAsBefore || isLatestDeleted; !shouldSnapshot {
		return latest, nil
	}

	snap := &snapshots.Snapshot{
		ID:          snapshotID,
		CommitSHA:   snapshotCommitSHA,
		CreatedAt:   time.Now(),
		WorkspaceID: workspaceID,
		CodebaseID:  codebaseID,
		Action:      action,
		DiffsCount:  &diffsCount,
	}

	if latest != nil {
		snap.PreviousSnapshotID = &latest.ID
	}

	if err := s.snapshotsRepo.Create(snap); err != nil {
		return nil, fmt.Errorf("can't create snapshot: %w", err)
	}

	if options.onView != nil || options.markAsLatestInWorkspace {
		isAuthoritativeView := ws.ViewID != nil && *ws.ViewID == *options.onView

		// If authoritative view, or explicitly asked to mark this as the latest snapshot
		if isAuthoritativeView || options.markAsLatestInWorkspace {
			ws.SetSnapshot(snap)
			if err := s.workspaceWriter.UpdateFields(ctx, ws.ID,
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

			// mark all snapshots as stale
			if err := s.statusesService.NotifyAllInWorkspace(ctx, ws.ID); err != nil {
				s.logger.Error("failed to notify statuses", zap.Error(err))
				// do not fail
			}
		}

		if isAuthoritativeView {
			if err := s.sendViewEvents(ctx, workspaceID, *options.onView); err != nil {
				return nil, fmt.Errorf("failed to send view events: %w", err)
			}
		}
	}

	logger.Info("snapshot created",
		zap.Duration("duration", time.Since(t0)))

	s.analyticsService.CaptureUser(ctx, ws.UserID, "created snapshot",
		analytics.CodebaseID(ws.CodebaseID),
		analytics.Property("diffs_count", diffsCount),
	)

	return snap, nil
}

func (s *Service) sendViewEvents(ctx context.Context, workspaceID, viewID string) error {
	// If this is a _suggestion_, send events to the view it's making suggestions to
	ws, err := s.workspaceReader.Get(workspaceID)
	if errors.Is(err, sql.ErrNoRows) {
		return nil
	} else if err != nil {
		return fmt.Errorf("could not get workspace: %w", err)
	}

	view, err := s.viewRepo.Get(viewID)
	if errors.Is(err, sql.ErrNoRows) {
		return nil
	} else if err != nil {
		return fmt.Errorf("could not get view: %w", err)
	}

	if ws.UserID == view.UserID {
		return nil
	}
	// find the owners views
	ownerViews, err := s.viewRepo.ListByCodebaseAndUser(ws.CodebaseID, ws.UserID)
	if errors.Is(err, sql.ErrNoRows) {
		return nil
	} else if err != nil {
		return fmt.Errorf("could not get workspace owner views: %w", err)
	}

	for _, ownerView := range ownerViews {
		if ownerView.WorkspaceID != workspaceID {
			continue
		}

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

func (s *Service) Diffs(ctx context.Context, snapshotID snapshots.ID, oo ...DiffsOption) ([]unidiff.FileDiff, error) {
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
	patchIDs *[]string
	withUser *users.User
}

type CopyOption func(*CopyOptions)

func CopyWithPatchIDs(patchIDs []string) CopyOption {
	return func(options *CopyOptions) {
		if options.patchIDs == nil {
			options.patchIDs = &patchIDs
		} else {
			*options.patchIDs = append(*options.patchIDs, patchIDs...)
		}
	}
}

func CopyWithUser(user *users.User) CopyOption {
	return func(options *CopyOptions) {
		options.withUser = user
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
func (s *Service) Copy(ctx context.Context, snapshotID snapshots.ID, oo ...CopyOption) (*snapshots.Snapshot, error) {
	snapshot, err := s.snapshotsRepo.Get(snapshotID)
	if err != nil {
		return nil, fmt.Errorf("failed to get snapshot: %w", err)
	}

	diffOptions := []DiffsOption{}
	options := getCopyOptions(oo...)
	if options.patchIDs != nil {
		diffOptions = append(diffOptions, DiffWithPatchIDs(*options.patchIDs))
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
		ID:          snapshots.ID(uuid.NewString()),
		CreatedAt:   time.Now(),
		CodebaseID:  snapshot.CodebaseID,
		WorkspaceID: snapshot.WorkspaceID,
		Action:      snapshot.Action,
		DiffsCount:  snapshot.DiffsCount,
	}

	gitSignature := git.Signature{
		Name:  "Sturdy",
		Email: "support@getsturdy.com",
		When:  time.Now(),
	}
	if options.withUser != nil {
		gitSignature.Name = options.withUser.Name
		gitSignature.Email = options.withUser.Email
	}

	if err := s.executorProvider.New().
		Write(vcs_view.CheckoutBranch(snapshot.WorkspaceID)).
		Write(func(repo vcs.RepoWriter) error {
			if err := repo.ApplyPatchesToWorkdir(patches); err != nil {
				return fmt.Errorf("failed to apply patches to workdir: %w", err)
			}

			commitID, err := vcs_snapshots.SnapshotOnViewRepo(ctx, s.logger, repo, newSnapshot.CodebaseID, newSnapshot.ID, gitSignature)
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
