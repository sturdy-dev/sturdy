package vcs

import (
	"errors"
	"fmt"
	"time"

	vcs_change "getsturdy.com/api/pkg/change/vcs"
	"getsturdy.com/api/pkg/unidiff"
	"getsturdy.com/api/vcs"
	"getsturdy.com/api/vcs/provider"

	git "github.com/libgit2/git2go/v33"
	"go.uber.org/zap"
)

type SnapshotOptions struct {
	patchIDsFilter *[]string
	revertCommitID *string
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

func WithRevertCommitID(commitID string) SnapshotOption {
	return func(opts *SnapshotOptions) {
		opts.revertCommitID = &commitID
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

func SnapshotOnViewRepo(logger *zap.Logger, repo vcs.RepoReaderGitWriter, codebaseID, snapshotID string, opts ...SnapshotOption) (string, error) {
	start := time.Now()

	options := snapshotOptions(opts...)

	if options.revertCommitID != nil {
		return "", errors.New("expected revertCommitID to be nul, was set")
	}
	if snapshotID == "" {
		return "", errors.New("snapshotID is not set")
	}

	preCommit, err := repo.HeadCommit()
	if err != nil {
		return "", fmt.Errorf("failed to find current head: %w", err)
	}
	defer preCommit.Free()

	patchIDs, err := snapshotPatchIDs(logger, repo, options)
	if err != nil {
		return "", err
	}

	var snapshotCommitID string

	// If no patches are specified, create a snapshot of the entire view ("git add -a")
	if len(patchIDs) == 0 {
		snapshotCommitID, err = repo.AddAndCommit(fmt.Sprintf("snapshot-%d", time.Now().Unix()))
	} else {
		sig := git.Signature{Name: "snapshot", Email: "snapshot@getsturdy.com", When: time.Now()}
		snapshotCommitID, err = vcs_change.CreateChangeFromPatchesOnRepo(logger, repo, codebaseID, patchIDs, fmt.Sprintf("snapshot-%d", time.Now().Unix()), sig)
	}
	if err != nil {
		return "", fmt.Errorf("failed to make snapshot: %w", err)
	}
	logger.Info("snapshot creation duration", zap.Duration("duration", time.Since(start)))

	err = repo.ResetMixed(preCommit.Id().String())
	if err != nil {
		return "", fmt.Errorf("failed to restore to workspace: %w", err)
	}

	// Push to upstream
	branchName := "snapshot-" + snapshotID
	if err := repo.CreateNewBranchAt(branchName, snapshotCommitID); err != nil {
		return "", fmt.Errorf("failed to create branch at snapshot: %w", err)
	}
	if err := repo.Push(logger, branchName); err != nil {
		return "", fmt.Errorf("failed to push snapshot branch: %w", err)
	}

	return snapshotCommitID, nil
}

func SnapshotOnTrunk(repo vcs.RepoGitWriter, workspaceID, snapshotID string, opts ...SnapshotOption) (string, error) {
	if snapshotID == "" {
		return "", errors.New("snapshotID is not set")
	}

	options := snapshotOptions(opts...)

	if options.revertCommitID == nil {
		return "", errors.New("expected revertCommitID to be set, got null")
	}

	workspaceHead, err := repo.BranchCommitID(workspaceID)
	if err != nil {
		return "", err
	}

	if err := repo.CreateNewBranchAt("snapshot-"+snapshotID, workspaceHead); err != nil {
		return "", err
	}

	newCommitID, err := repo.RevertOnBranch(*options.revertCommitID, "snapshot-"+snapshotID)
	if err != nil {
		return "", err
	}

	return newCommitID, nil
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

func Restore(logger *zap.Logger, viewProvider provider.ViewProvider, codebaseID, workspaceID, viewID, snapshotID, snapshotCommitID string) error {
	repo, err := viewProvider.ViewRepo(codebaseID, viewID)
	if err != nil {
		return err
	}
	return RestoreRepo(logger, repo, codebaseID, workspaceID, snapshotID, snapshotCommitID)
}

func RestoreRepo(logger *zap.Logger, repo vcs.RepoWriter, codebaseID, workspaceID, snapshotID, snapshotCommitID string) error {
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
		// return fmt.Errorf("failed to pull large files: %w", err)
	}
	logger.Info("snapshot restore large files pulled", zap.Duration("duration", time.Since(t0)))

	// Mixed reset to the snapshot commits parent (a user authored commit)
	if err := repo.ResetMixed(parents[0]); err != nil {
		return fmt.Errorf("failed to reset mixed: %w", err)
	}

	if err := repo.ForcePush(logger, workspaceID); err != nil {
		return fmt.Errorf("failed to push: %w", err)
	}

	return nil
}

func Diff(logger *zap.Logger, viewProvider provider.ViewProvider, codebaseID, viewID, snapshotCommitID, parentSnapshotCommitID string) ([]unidiff.FileDiff, error) {
	repo, err := viewProvider.ViewRepo(codebaseID, viewID)
	if err != nil {
		return nil, err
	}

	diffs, err := repo.DiffCommits(snapshotCommitID, parentSnapshotCommitID)
	if err != nil {
		return nil, err
	}
	defer diffs.Free()

	res, err := unidiff.NewUnidiff(unidiff.NewGitPatchReader(diffs), logger).Decorate()
	if err != nil {
		return nil, err
	}

	return res, nil
}
