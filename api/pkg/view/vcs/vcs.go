package vcs

import (
	"fmt"

	"getsturdy.com/api/pkg/codebases"
	"getsturdy.com/api/pkg/snapshots"
	"getsturdy.com/api/vcs"
	"getsturdy.com/api/vcs/provider"
)

func Create(repoProvider provider.RepoProvider, codebaseID codebases.ID, checkoutBranchName, viewID string) error {
	view, err := vcs.CloneRepo(repoProvider.TrunkPath(codebaseID), repoProvider.ViewPath(codebaseID, viewID))
	if err != nil {
		return fmt.Errorf("failed to create a view of %s: %w", codebaseID, err)
	}
	return checkoutBranch(view, checkoutBranchName)
}

func SetWorkspace(viewProvider provider.ViewProvider, codebaseID codebases.ID, viewID, workspaceID string) error {
	repo, err := viewProvider.ViewRepo(codebaseID, viewID)
	if err != nil {
		return fmt.Errorf("failed find codebaseID %s: %w", codebaseID, err)
	}
	return checkoutBranch(repo, workspaceID)
}

func CheckoutBranch(branchName string) func(vcs.RepoWriter) error {
	return func(repo vcs.RepoWriter) error {
		return checkoutBranch(repo, branchName)
	}
}

func checkoutBranch(repo vcs.RepoWriter, branchName string) error {
	headBranch, err := repo.HeadBranch()
	if err != nil {
		return fmt.Errorf("failed to get head branch: %w", err)
	}
	if headBranch == branchName {
		return nil
	}
	if err := repo.FetchBranch(branchName); err != nil {
		return fmt.Errorf("failed to fetch branch '%s': %w", branchName, err)
	}
	if err := repo.CreateBranchTrackingUpstream(branchName); err != nil {
		return fmt.Errorf("failed to create branch '%s': %w", branchName, err)
	}
	if err := repo.CheckoutBranchWithForce(branchName); err != nil {
		return fmt.Errorf("failed to checkout branch '%s': %w", branchName, err)
	}
	if err := repo.CleanStaged(); err != nil {
		return fmt.Errorf("failed to clean index after checkout '%s': %w", branchName, err)
	}
	return nil
}

func CheckoutSnapshot(snapshot *snapshots.Snapshot) func(vcs.RepoWriter) error {
	return func(repo vcs.RepoWriter) error {
		copyBranchName := snapshot.BranchName()
		copyParentCommitsIDs, err := repo.GetCommitParents(snapshot.CommitSHA)
		if err != nil {
			return fmt.Errorf("failed to get commit parents: %w", err)
		}
		if len(copyParentCommitsIDs) != 1 {
			return fmt.Errorf("unexpected number of parents=%d", len(copyParentCommitsIDs))
		}
		preCommitID := copyParentCommitsIDs[0]
		if err := repo.CreateBranchTrackingUpstream(copyBranchName); err != nil {
			return fmt.Errorf("failed to create branch on target: %w", err)
		}
		if err := repo.CheckoutBranchWithForce(copyBranchName); err != nil {
			return fmt.Errorf("failed to checkout branch on target: %w", err)
		}
		if err := repo.CreateBranchTrackingUpstream(snapshot.WorkspaceID); err != nil {
			return fmt.Errorf("failed to create workspace branch on target: %w", err)
		}
		if err := repo.ResetMixed(preCommitID); err != nil {
			return fmt.Errorf("failed to restore to parent on target: %w", err)
		}
		if err := repo.CheckoutBranchSafely(snapshot.WorkspaceID); err != nil {
			return fmt.Errorf("failed to checkout workspace on target: %w", err)
		}
		return nil
	}
}
