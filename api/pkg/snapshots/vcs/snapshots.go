package vcs

import (
	"errors"
	"fmt"
	"time"

	vcs_change "getsturdy.com/api/pkg/changes/vcs"
	"getsturdy.com/api/pkg/codebases"
	"getsturdy.com/api/pkg/snapshots"
	"getsturdy.com/api/pkg/unidiff"
	"getsturdy.com/api/vcs"

	git "github.com/libgit2/git2go/v33"
	"go.uber.org/zap"
)

type SnapshotOptions struct {
	patchIDsFilter       *[]string
	revertCommitHeadBase *[2]*string
	commitMessage        string
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

func WithRevert(head string, base *string) SnapshotOption {
	return func(opts *SnapshotOptions) {
		opts.revertCommitHeadBase = &[2]*string{&head, base}
	}
}

func WithCommitMessage(msg string) SnapshotOption {
	return func(opts *SnapshotOptions) {
		opts.commitMessage = msg
	}
}

func snapshotOptions(opts ...SnapshotOption) *SnapshotOptions {
	options := &SnapshotOptions{}
	for _, applyOption := range opts {
		applyOption(options)
	}
	return options
}

func snapshotPatchIDs(logger *zap.Logger, repo vcs.RepoGitReader, options *SnapshotOptions) ([]string, error) {
	if options.patchIDsFilter != nil {
		return *options.patchIDsFilter, nil
	}
	return allPatchIDs(logger, repo)
}

func allPatchIDs(logger *zap.Logger, repo vcs.RepoGitReader) ([]string, error) {
	diffs, err := repo.CurrentDiffNoIndex()
	if err != nil {
		return nil, fmt.Errorf("failed to get current diff: %w", err)
	}
	defer diffs.Free()

	differ := unidiff.NewUnidiff(unidiff.NewGitPatchReader(diffs), logger).WithExpandedHunks()
	fileDiffs, err := differ.Decorate()
	if err != nil {
		return nil, fmt.Errorf("failed to build diffs: %w", err)
	}

	var patchIDs []string
	for _, d := range fileDiffs {
		for _, h := range d.Hunks {
			patchIDs = append(patchIDs, h.ID)
		}
	}

	return patchIDs, nil
}

func SnapshotOnViewRepo(logger *zap.Logger, repo vcs.RepoReaderGitWriter, codebaseID codebases.ID, snapshotID string, signature git.Signature, opts ...SnapshotOption) (string, error) {
	start := time.Now()

	options := snapshotOptions(opts...)

	if snapshotID == "" {
		return "", errors.New("snapshotID is not set")
	}
	if options.revertCommitHeadBase != nil {
		return "", errors.New("expected revertCommitID to be nul, was set")
	}

	decoratedLogger := logger.With(
		zap.String("codebase_id", codebaseID.String()),
		zap.String("snapshot_id", snapshotID),
		zap.String("repo_path", repo.Path()),
	)

	preCommit, err := repo.HeadCommit()
	if err != nil {
		return "", fmt.Errorf("failed to find current head: %w", err)
	}
	defer preCommit.Free()

	patchIDs, err := snapshotPatchIDs(decoratedLogger, repo, options)
	if err != nil {
		return "", err
	}

	var snapshotCommitID string

	decoratedLogger.Info("creating snapshot", zap.Int("patch_count", len(patchIDs)))

	commitMessage := fmt.Sprintf("snapshot-%d", time.Now().Unix())
	if options.commitMessage != "" {
		commitMessage = options.commitMessage
	}

	// If no patches are specified, create a snapshot of the entire view ("git add -a")
	if len(patchIDs) == 0 {
		snapshotCommitID, err = repo.AddAndCommitWithSignature(commitMessage, signature)
	} else {
		snapshotCommitID, err = vcs_change.CreateChangeFromPatchesOnRepo(decoratedLogger, repo, codebaseID, patchIDs, commitMessage, signature)
	}
	if err != nil {
		return "", fmt.Errorf("failed to make snapshot: %w", err)
	}
	decoratedLogger.Info("snapshot creation duration", zap.Duration("duration", time.Since(start)))

	err = repo.ResetMixed(preCommit.Id().String())
	if err != nil {
		return "", fmt.Errorf("failed to restore to workspace: %w", err)
	}

	// Push to upstream
	branchName := "snapshot-" + snapshotID
	if err := repo.CreateNewBranchAt(branchName, snapshotCommitID); err != nil {
		return "", fmt.Errorf("failed to create branch at snapshot branchName=%s snapshotCommitID=%s: %w", branchName, snapshotCommitID, err)
	}
	if err := repo.Push(decoratedLogger, branchName); err != nil {
		return "", fmt.Errorf("failed to push snapshot branch branchName=%s snapshotCommitID=%s: %w", branchName, snapshotCommitID, err)
	}

	return snapshotCommitID, nil
}

func SnapshotOnViewRepoWithRevert(repo vcs.RepoWriter, logger *zap.Logger, snapshotID string, opts ...SnapshotOption) (string, error) {
	if snapshotID == "" {
		return "", errors.New("snapshotID is not set")
	}

	options := snapshotOptions(opts...)

	if options.revertCommitHeadBase == nil {
		return "", errors.New("expected revertCommitID to be set, got null")
	}

	hb := *options.revertCommitHeadBase

	// hb[0] is the "head" the change that's going to be the new base of the workspace
	// hb[1] is the "base" (the parent change of the "head") and is optional. If it's nil, diff against the root of the codebase.
	if hb[0] == nil {
		return "", errors.New("revertCommitHeadBase[0] (the head) is nil")
	}

	head := *hb[0]
	base := hb[1]

	if err := repo.CreateAndCheckoutBranchAtCommit(head, "snapshot-pre-revert-"+snapshotID); err != nil {
		return "", fmt.Errorf("failed to create pre-revert branch: %w", err)
	}

	var diff *git.Diff
	var err error
	if base == nil {
		diff, err = repo.DiffRootToCommit(head)
	} else {
		diff, err = repo.DiffCommits(head, *base)
	}
	if err != nil {
		return "", fmt.Errorf("failed to create reversed diff: %w", err)
	}

	patches, err := unidiff.NewUnidiff(unidiff.NewGitPatchReader(diff), logger).PatchesBytes()
	if err != nil {
		return "", fmt.Errorf("failed to get patches: %w", err)
	}

	if err := repo.ApplyPatchesToWorkdir(patches); err != nil {
		return "", fmt.Errorf("failed to apply patches: %w", err)
	}

	snapshotCommitID, err := repo.AddAndCommit(fmt.Sprintf("snapshot-%d", time.Now().Unix()))
	if err != nil {
		return "", fmt.Errorf("failed to commit snapshot: %w", err)
	}

	err = repo.ResetMixed(*hb[0])
	if err != nil {
		return "", fmt.Errorf("failed to restore to workspace: %w", err)
	}

	branchName := "snapshot-" + snapshotID

	// Push to upstream
	if err := repo.CreateNewBranchAt(branchName, snapshotCommitID); err != nil {
		return "", fmt.Errorf("failed to create branch at snapshot branchName=%s snapshotCommitID=%s: %w", branchName, snapshotCommitID, err)
	}
	if err := repo.Push(logger, branchName); err != nil {
		return "", fmt.Errorf("failed to push snapshot branch branchName=%s snapshotCommitID=%s: %w", branchName, snapshotCommitID, err)
	}

	return snapshotCommitID, nil
}

func SnapshotOnExistingCommit(repo vcs.RepoGitWriter, snapshotID, existingCommitID string) (string, error) {
	if snapshotID == "" {
		return "", errors.New("snapshotID is not set")
	}
	if err := repo.CreateNewBranchAt("snapshot-"+snapshotID, existingCommitID); err != nil {
		return "", err
	}
	return existingCommitID, nil
}

func Restore(logger *zap.Logger, snapshot *snapshots.Snapshot) func(vcs.RepoWriter) error {
	return func(repo vcs.RepoWriter) error {
		return RestoreRepo(logger, repo, snapshot.ID, snapshot.CommitSHA)
	}
}

func RestoreRepo(logger *zap.Logger, repo vcs.RepoWriter, snapshotID, snapshotCommitID string) error {
	if err := repo.FetchBranch("snapshot-" + snapshotID); err != nil {
		return fmt.Errorf("failed to fetch: %w", err)
	}

	parents, err := repo.GetCommitParents(snapshotCommitID)
	if err != nil {
		return fmt.Errorf("failed to get parents: %w", err)
	}
	if len(parents) != 1 {
		return fmt.Errorf("unexpected number of parents: %d", len(parents))
	}

	// Reset HARD to the snapshot commit
	if err := repo.ResetHard(snapshotCommitID); err != nil {
		return fmt.Errorf("failed to reset: %w", err)
	}

	t0 := time.Now()
	if err := repo.LargeFilesPull(); err != nil {
		logger.Warn("failed to pull large files", zap.Error(err))
	}
	logger.Info("snapshot restore large files pulled", zap.Duration("duration", time.Since(t0)))

	// Mixed reset to the snapshot commits parent
	if err := repo.ResetMixed(parents[0]); err != nil {
		return fmt.Errorf("failed to reset mixed: %w", err)
	}

	return nil
}
