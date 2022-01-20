package vcs

import (
	"fmt"

	"github.com/google/uuid"

	vcsvcs "getsturdy.com/api/vcs"
)

const UnsavedCommitMessage = "Unsaved workspace changes"

func SyncSingleCommitOnBranch(repo vcsvcs.RepoWriter, syncCommit, remoteName, branchName string) error {
	err := repo.FetchBranch(branchName)
	if err != nil {
		return fmt.Errorf("fetch origin failed: %w", err)
	}

	onto, err := repo.RemoteBranchCommit(remoteName, branchName)
	if err != nil {
		return err
	}
	defer onto.Free()

	syncingBranchName := fmt.Sprintf("syncing-%s", uuid.NewString())

	if err := repo.CreateAndCheckoutBranchAtCommit(onto.Id().String(), syncingBranchName); err != nil {
		return fmt.Errorf("create and checkout branch failed: %w", err)
	}

	_, conflicted, _, err := repo.CherryPickOnto(syncCommit, onto.Id().String())
	if conflicted {
		return fmt.Errorf("could not sync, had conflicts")
	}
	if err != nil {
		return fmt.Errorf("cherry pick failed: %w", err)
	}

	if err := repo.MoveBranchToHEAD(branchName); err != nil {
		return err
	}

	return nil
}
