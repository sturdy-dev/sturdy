package vcs

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path"

	"getsturdy.com/api/pkg/codebases"
	vcs_sync "getsturdy.com/api/pkg/sync/vcs"
	"getsturdy.com/api/pkg/unidiff"
	"getsturdy.com/api/vcs"
	"getsturdy.com/api/vcs/executor"

	"github.com/google/uuid"
	"go.uber.org/zap"

	git "github.com/libgit2/git2go/v33"
)

func CreateAndLandFromView(
	viewRepo vcs.RepoWriter,
	logger *zap.Logger,
	codebaseID codebases.ID,
	workspaceID string,
	patchIDs []string,
	message string,
	signature git.Signature,
	diffOpts ...vcs.DiffOption,
) (string, func(vcs.RepoGitWriter) error, error) {
	preCreateBranchHead, err := viewRepo.BranchCommitID(workspaceID)
	if err != nil {
		return "", nil, fmt.Errorf("failed to get branch head: %w", err)
	}
	creationBranchName := fmt.Sprintf("create-change-%s", uuid.NewString())

	logger = logger.With(zap.String("creation_branch_name", creationBranchName))

	if err := viewRepo.CreateNewBranchAt(creationBranchName, preCreateBranchHead); err != nil {
		return "", nil, fmt.Errorf("failed to create new branch to use during creation: %w", err)
	}
	if err := viewRepo.CheckoutBranchSafely(creationBranchName); err != nil {
		return "", nil, fmt.Errorf("failed to checkout the new branch: %w", err)
	}

	defer func() {
		if err == nil {
			return
		}
		// if something went wrong, restore the view to how it was
		logger.Error("create and land failed, trying to restore", zap.Error(err))

		if err := viewRepo.CheckoutBranchWithForce(creationBranchName); err != nil {
			logger.Error("failed to checkout the creation branch", zap.Error(err))
		}
		if err := viewRepo.ResetMixed(preCreateBranchHead); err != nil {
			logger.Error("failed to reset to the pre create head", zap.Error(err))
		}
		if err := viewRepo.MoveBranchToHEAD(workspaceID); err != nil {
			logger.Error("failed to move the head", zap.Error(err))
		}
		if err := viewRepo.CheckoutBranchSafely(workspaceID); err != nil {
			logger.Error("failed to checkout the workspace branch", zap.Error(err))
		}
		if err := viewRepo.LargeFilesPull(); err != nil {
			// Log and continue (repo can have LFS files from outside of Sturdy)
			logger.Warn("failed to pull large files", zap.Error(err))
		}

		logger.Info("successfully restored view after failed landing")
	}()

	createdCommitID, err := CreateChangeFromPatchesOnRepo(logger, viewRepo, codebaseID, patchIDs, message, signature, diffOpts...)
	if err != nil {
		return "", nil, fmt.Errorf("failed to create the new change: %w", err)
	}

	if err = vcs_sync.FastLand(viewRepo, createdCommitID); err != nil {
		return "", nil, fmt.Errorf("landing failed: %w", err)
	}

	// move the workspace branch to be the same as the new sturdytrunk
	if err := viewRepo.MoveBranch(workspaceID, "sturdytrunk"); err != nil {
		return "", nil, fmt.Errorf("failed to move workspace to new trunk: %w", err)
	}

	if err := viewRepo.CheckoutBranchWithForce(workspaceID); err != nil {
		return "", nil, fmt.Errorf("failed to checkout workspace branch: %w", err)
	}

	// LFS Pull
	if err := viewRepo.LargeFilesPull(); err != nil {
		// Log and continue (repo can have LFS files from outside of Sturdy)
		logger.Warn("failed to pull large files", zap.Error(err))
	}

	newBranchCommit, err := viewRepo.BranchCommitID(workspaceID)
	if err != nil {
		return "", nil, fmt.Errorf("failed to get commit after land: %w", err)
	}

	// will be executed once the new state has been recorded in the databases
	pushFunc := func(viewRepo vcs.RepoGitWriter) error {
		if err := viewRepo.Push(logger, "sturdytrunk"); err != nil {
			return fmt.Errorf("push failed: %w", err)
		}

		if err := viewRepo.ForcePush(logger, workspaceID); err != nil {
			return fmt.Errorf("failed to push updated workspace branch: %w", err)
		}

		return nil
	}

	return newBranchCommit, pushFunc, nil
}

func CreateChangeFromPatchesOnRepo(logger *zap.Logger, r vcs.RepoReaderGitWriter, codebaseID codebases.ID, patchIDs []string, message string, signature git.Signature, diffOpts ...vcs.DiffOption) (string, error) {
	treeID, err := CreateChangesTreeFromPatches(logger, r, codebaseID, patchIDs, diffOpts...)
	if err != nil {
		return "", err
	}

	// No changes where added
	if treeID == nil {
		return "", fmt.Errorf("no changes to add")
	}

	changeID, err := r.CommitIndexTree(treeID, message, signature)
	if err != nil {
		return "", fmt.Errorf("failed save change: %w", err)
	}

	return changeID, nil
}

// CreateChangesTreeFromPatches creates a git-tree based on the inputs.
// If patchIDs is non-nil, the slice will be passed as a filter to unidiff.WithHunksFilter
func CreateChangesTreeFromPatches(logger *zap.Logger, r vcs.RepoReaderGitWriter, codebaseID codebases.ID, patchIDs []string, diffOpts ...vcs.DiffOption) (*git.Oid, error) {
	err := r.CleanStaged()
	if err != nil {
		return nil, fmt.Errorf("failed to clean staged codebase=%s %w", codebaseID, err)
	}

	currentDiff, err := r.CurrentDiff(diffOpts...)
	if err != nil {
		return nil, err
	}
	defer currentDiff.Free()

	diffs := unidiff.NewUnidiff(unidiff.NewGitPatchReader(currentDiff), logger).
		// Expand and filter hunks
		WithExpandedHunks()

	if patchIDs != nil {
		diffs = diffs.WithHunksFilter(patchIDs...)
	}

	// Join the hunks again. This is needed to be able to apply hunks where all of them have a file move, etc.
	binaryDiffs, nonBinaryDiffs, err := diffs.WithJoiner().DecorateSeparateBinary()
	if err != nil {
		return nil, err
	}

	var nonBinaryPatches [][]byte
	for _, diff := range nonBinaryDiffs {
		for _, hunk := range diff.Hunks {
			nonBinaryPatches = append(nonBinaryPatches, []byte(hunk.Patch))
		}
	}

	var treeID *git.Oid
	if len(nonBinaryPatches) > 0 {
		treeID, err = r.ApplyPatchesToIndex(nonBinaryPatches)
		if err != nil {
			for id, patch := range nonBinaryPatches {
				logger.Warn("failed to add patches", zap.Int("key", id), zap.String("patch", string(patch)))
			}
			return nil, fmt.Errorf("failed to add patches: %w", err)
		}
	}

	var binaryFiles []string
	for _, binDiff := range binaryDiffs {
		if binDiff.IsNew {
			binaryFiles = append(binaryFiles, binDiff.NewName)
		} else if binDiff.IsDeleted {
			binaryFiles = append(binaryFiles, binDiff.OrigName)
		} else if binDiff.NewName != binDiff.OrigName {
			binaryFiles = append(binaryFiles,
				binDiff.OrigName,
				binDiff.NewName,
			)
		} else {
			// Modified without rename
			binaryFiles = append(binaryFiles, binDiff.NewName)
		}
	}

	// Separate large (tracked in LFS) and small files (tracked in Git)
	var largeFiles []string
	var smallFiles []string
	for _, fName := range binaryFiles {
		fullPath := path.Join(r.Path(), fName)
		if stat, err := os.Stat(fullPath); err == nil && stat.Size() > 1_000_000 {
			largeFiles = append(largeFiles, fName)
		} else {
			smallFiles = append(smallFiles, fName)
		}
	}

	if len(largeFiles) > 0 {
		patches, err := r.LargeFilesClean(codebaseID, largeFiles)
		if err != nil {
			return nil, fmt.Errorf("failed to clean large files: %w", err)
		}

		treeID, err = r.ApplyPatchesToIndex(patches)
		if err != nil {
			return nil, fmt.Errorf("failed to add lfs patches: %w", err)
		}
	}

	if len(smallFiles) > 0 {
		treeID, err = r.AddFilesToIndex(smallFiles)
		if err != nil {
			return nil, fmt.Errorf("failed to add binary files: %w", err)
		}
	}

	return treeID, nil
}

func AddToGitignore(executorProvider executor.Provider, codebaseID codebases.ID, viewID, ignorePath string) error {
	executor := executorProvider.New().Read(func(repo vcs.RepoReader) error {
		ignoreFilePath := path.Join(repo.Path(), ".gitignore")

		all, err := ioutil.ReadFile(ignoreFilePath)
		if errors.Is(err, os.ErrNotExist) {
			all = []byte{} // Create new file
		} else if err != nil {
			return err
		}

		if len(all) > 0 {
			all = bytes.TrimRight(all, "\n\r")
			all = append(all, []byte("\n")...)
		}
		all = append(all, []byte(ignorePath+"\n")...)

		err = ioutil.WriteFile(ignoreFilePath, all, 0o644)
		if err != nil {
			return err
		}
		return nil
	})

	if err := executor.ExecView(codebaseID, viewID, "AddToGitignore"); err != nil {
		return err
	}

	return nil
}
