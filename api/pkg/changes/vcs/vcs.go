package vcs

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path"

	"getsturdy.com/api/pkg/codebases"
	"getsturdy.com/api/pkg/unidiff"
	"getsturdy.com/api/vcs"
	"getsturdy.com/api/vcs/executor"

	"go.uber.org/zap"

	git "github.com/libgit2/git2go/v33"
)

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
