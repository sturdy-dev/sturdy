package vcs

import (
	"fmt"

	"mash/pkg/unidiff"
	"mash/pkg/unidiff/lfs"
	"mash/vcs"

	"go.uber.org/zap"
)

func RemoveWithPatches(logger *zap.Logger, patches [][]byte, patchIDs ...string) func(vcs.RepoWriter) error {
	return func(r vcs.RepoWriter) error {
		return remove(r, logger, unidiff.NewBytesPatchReader(patches), patchIDs...)
	}
}

func Remove(logger *zap.Logger, patchIDs ...string) func(vcs.RepoWriter) error {
	return func(r vcs.RepoWriter) error {
		currentDiff, err := r.CurrentDiff()
		if err != nil {
			return err
		}
		defer currentDiff.Free()

		return remove(r, logger, unidiff.NewGitPatchReader(currentDiff), patchIDs...)
	}
}

func remove(r vcs.RepoWriter, logger *zap.Logger, patchReader unidiff.PatchReader, patchIDs ...string) error {
	lfsFilter, err := lfs.NewIgnoreLfsSmudgedFilter(r)
	if err != nil {
		return err
	}

	binaryDiffs, nonBinaryDiffs, err := unidiff.NewUnidiff(patchReader, logger).
		WithExpandedHunks().
		WithHunksFilter(patchIDs...).
		WithInverter().
		WithFilterFunc(lfsFilter).
		DecorateSeparateBinary()
	if err != nil {
		return err
	}

	var nonBinaryPatches [][]byte
	for _, diff := range nonBinaryDiffs {
		if diff.IsNew || diff.IsDeleted {
			binaryDiffs = append(binaryDiffs, diff)
			continue
		}

		for _, hunk := range diff.Hunks {
			nonBinaryPatches = append(nonBinaryPatches, []byte(hunk.Patch))
		}
	}

	if err := r.ApplyPatchesToWorkdir(nonBinaryPatches); err != nil {
		return fmt.Errorf("failed to add inverted patches: %w", err)
	}

	for _, binDiff := range binaryDiffs {
		// The status has been inverted from what the user is experiencing.
		switch {
		case binDiff.IsDeleted:
			if err := r.CheckoutFile(binDiff.OrigName); err != nil {
				return fmt.Errorf("failed to reset binary file (was deleted): %w", err)
			}

		case binDiff.IsNew:
			if err := r.DeleteFile(binDiff.NewName); err != nil {
				return fmt.Errorf("failed to delete binary file (was new): %w", err)
			}

		case !binDiff.IsNew && !binDiff.IsDeleted && binDiff.OrigName == binDiff.NewName: // modified in-place
			if err := r.CheckoutFile(binDiff.OrigName); err != nil {
				return fmt.Errorf("failed to reset binary file (was simple modify): %w", err)
			}

		case binDiff.OrigName != binDiff.NewName: // moved
			if err := r.CheckoutFile(binDiff.OrigName); err != nil {
				return fmt.Errorf("undo renamed: failed to checkout original version: %w", err)
			}
			if err := r.DeleteFile(binDiff.NewName); err != nil {
				return fmt.Errorf("undo renamed: failed to delete new version: %w", err)
			}

		default:
			return fmt.Errorf("did not know how to revert this file")
		}

	}

	if err := r.LargeFilesPull(); err != nil {
		return fmt.Errorf("failed to pull large files after revert: %w", err)
	}

	return nil

}
