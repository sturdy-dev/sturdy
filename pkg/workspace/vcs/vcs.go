package vcs

import (
	"fmt"

	"mash/vcs"
)

func Create(repo vcs.Repo, workspaceID string) error {
	if err := repo.CreateNewBranchOnHEAD(workspaceID); err != nil {
		return fmt.Errorf("failed to create workspaceID %s: %w", workspaceID, err)
	}
	return nil
}

func CreateAtChange(repo vcs.Repo, workspaceID, changeID string) error {
	if err := repo.CreateNewBranchAt(workspaceID, changeID); err != nil {
		return fmt.Errorf("failed to create workspaceID %s: %w", workspaceID, err)
	}
	return nil
}
