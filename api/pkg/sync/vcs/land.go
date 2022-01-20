package vcs

import (
	"fmt"

	"getsturdy.com/api/vcs"
)

func FastLand(viewRepo vcs.RepoWriter, commitID string) (err error) {
	if err = viewRepo.FetchBranch("sturdytrunk"); err != nil {
		return fmt.Errorf("failed to fetch before fastland: %w", err)
	}

	if err := SyncSingleCommitOnBranch(viewRepo, commitID, "origin", "sturdytrunk"); err != nil {
		return fmt.Errorf("failed to land: %w", err)
	}

	return nil
}
