package vcs

import (
	"fmt"

	"getsturdy.com/api/vcs"
	"getsturdy.com/api/vcs/provider"
)

func Create(repoProvider provider.RepoProvider, codebaseID, checkoutBranchName, viewID string) error {
	view, err := vcs.CloneRepo(repoProvider.TrunkPath(codebaseID), repoProvider.ViewPath(codebaseID, viewID))
	if err != nil {
		return fmt.Errorf("failed to create a view of %s: %w", codebaseID, err)
	}

	if err := view.FetchBranch(checkoutBranchName); err != nil {
		return fmt.Errorf("failed to fetch branch: %w", err)
	}

	if err := view.CreateBranchTrackingUpstream(checkoutBranchName); err != nil {
		return fmt.Errorf("failed to create branch: %w", err)
	}

	if err := view.CheckoutBranchWithForce(checkoutBranchName); err != nil {
		return fmt.Errorf("failed to checkout branch: %w", err)
	}

	return nil
}

func SetWorkspace(viewProvider provider.ViewProvider, codebaseID, viewID, workspaceID string) error {
	repo, err := viewProvider.ViewRepo(codebaseID, viewID)
	if err != nil {
		return fmt.Errorf("failed find codebaseID %s: %w", codebaseID, err)
	}
	return SetWorkspaceRepo(repo, workspaceID)
}

func SetWorkspaceRepo(repo vcs.RepoWriter, workspaceID string) error {
	err := repo.FetchBranch(workspaceID)
	if err != nil {
		return fmt.Errorf("failed to update view (step 1): %w", err)
	}

	err = repo.CreateBranchTrackingUpstream(workspaceID)
	if err != nil {
		return fmt.Errorf("failed to update view (step 2): %w", err)
	}

	err = repo.CheckoutBranchWithForce(workspaceID)
	if err != nil {
		return fmt.Errorf("failed to update view (step 3): %w", err)
	}

	return nil
}
