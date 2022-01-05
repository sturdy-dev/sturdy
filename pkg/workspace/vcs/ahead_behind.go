package vcs

import (
	"fmt"

	"mash/vcs"
	"mash/vcs/provider"
)

func UpToDateWithTrunk(repo vcs.Repo, workspaceID string) (bool, error) {
	trunkHEAD, err := repo.BranchCommitID("sturdytrunk")
	if err != nil {
		// If sturdytrunk doesn't exist (such as when an empty repository has been imported), treat it as up to date
		return true, nil
	}
	return repo.BranchHasCommit(workspaceID, trunkHEAD)
}

func BehindAheadCount(trunkProvider provider.TrunkProvider, codebaseID, workspaceID string) (behind, ahead int, err error) {
	repo, err := trunkProvider.TrunkRepo(codebaseID)
	if err != nil {
		return 0, 0, fmt.Errorf("failed find codebaseID %s: %w", codebaseID, err)
	}
	defer repo.Free()

	trunkHEAD, err := repo.BranchCommitID("sturdytrunk")
	if err != nil {
		return 0, 0, err
	}

	workspaceHEAD, err := repo.BranchCommitID(workspaceID)
	if err != nil {
		return 0, 0, err
	}

	behindCommits, err := repo.RevlistCherryPickRightOnly(workspaceHEAD, trunkHEAD)
	if err != nil {
		return 0, 0, fmt.Errorf("failed to list commits behind: %w", err)
	}

	aheadCommits, err := repo.RevlistCherryPickLeftOnly(workspaceHEAD, trunkHEAD)
	if err != nil {
		return 0, 0, fmt.Errorf("failed to list commits ahead: %w", err)
	}

	return len(behindCommits), len(aheadCommits), nil
}
