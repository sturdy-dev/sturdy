package vcs

import (
	"getsturdy.com/api/vcs"
)

func UpToDateWithTrunk(repo vcs.RepoGitReader, workspaceID string) (bool, error) {
	trunkHEAD, err := repo.BranchCommitID("sturdytrunk")
	if err != nil {
		// If sturdytrunk doesn't exist (such as when an empty repository has been imported), treat it as up to date
		return true, nil
	}
	return repo.BranchHasCommit(workspaceID, trunkHEAD)
}
