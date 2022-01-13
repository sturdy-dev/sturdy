package vcs

import (
	"mash/vcs"
)

func UpToDateWithTrunk(repo vcs.Repo, workspaceID string) (bool, error) {
	trunkHEAD, err := repo.BranchCommitID("sturdytrunk")
	if err != nil {
		// If sturdytrunk doesn't exist (such as when an empty repository has been imported), treat it as up to date
		return true, nil
	}
	return repo.BranchHasCommit(workspaceID, trunkHEAD)
}
