package vcs

import (
	"fmt"
	"mash/vcs"
	"mash/vcs/provider"
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

func CopyUpload(logger *zap.Logger, viewProvider provider.ViewProvider, codebaseID, sourceViewID string) (preCommitID, copyBranchName string, err error) {
	sourceRepo, err := viewProvider.ViewRepo(codebaseID, sourceViewID)
	if err != nil {
		return "", "", err
	}
	defer sourceRepo.Free()

	preCommit, err := sourceRepo.HeadCommit()
	if err != nil {
		return "", "", fmt.Errorf("failed to find current head: %w", err)
	}
	defer preCommit.Free()

	snapshotCommitID, err := sourceRepo.AddAndCommit(fmt.Sprintf("copy-%d", time.Now().Unix()))
	if err != nil {
		return "", "", fmt.Errorf("failed to add snapshot: %w", err)
	}
	err = sourceRepo.ResetMixed(preCommit.Id().String())
	if err != nil {
		return "", "", fmt.Errorf("failed to restore to workspace: %w", err)
	}

	copyBranchName = "copy-" + uuid.New().String()
	err = sourceRepo.CreateNewBranchAt(copyBranchName, snapshotCommitID)
	if err != nil {
		return "", "", fmt.Errorf("failed to create copy branch: %w", err)
	}

	err = sourceRepo.Push(logger, copyBranchName)
	if err != nil {
		return "", "", fmt.Errorf("failed to push copy branch: %w", err)
	}

	return preCommit.Id().String(), copyBranchName, nil
}

func Copy(logger *zap.Logger, viewProvider provider.ViewProvider, codebaseID, sourceWorkspaceID, sourceViewID, targetViewID string) (string, error) {
	preCommitID, copyBranchName, err := CopyUpload(logger, viewProvider, codebaseID, sourceViewID)
	if err != nil {
		return "", err
	}

	targetRepo, err := viewProvider.ViewRepo(codebaseID, targetViewID)
	if err != nil {
		return "", err
	}
	defer targetRepo.Free()

	return checkoutCopy(targetRepo, copyBranchName, sourceWorkspaceID, preCommitID)
}

func CheckoutWorkspaceSnapshot(repoProvider provider.RepoProvider, codebaseID, sourceWorkspaceID, workspaceSnapshotBranchName, targetViewID string) (string, error) {
	copyBranchName := workspaceSnapshotBranchName

	trunkRepo, err := repoProvider.TrunkRepo(codebaseID)
	if err != nil {
		return "", err
	}
	defer trunkRepo.Free()

	copyCommitID, err := trunkRepo.BranchCommitID(copyBranchName)
	if err != nil {
		return "", fmt.Errorf("failed to get branch HEAD commit: %w", err)
	}
	copyParentCommitsIDs, err := trunkRepo.GetCommitParents(copyCommitID)
	if err != nil {
		return "", fmt.Errorf("failed to get commit parents: %w", err)
	}
	if len(copyParentCommitsIDs) != 1 {
		return "", fmt.Errorf("unexpected number of parents=%d", len(copyParentCommitsIDs))
	}
	preCommitID := copyParentCommitsIDs[0]

	targetRepo, err := repoProvider.ViewRepo(codebaseID, targetViewID)
	if err != nil {
		return "", err
	}
	defer targetRepo.Free()

	return checkoutCopy(targetRepo, copyBranchName, sourceWorkspaceID, preCommitID)
}

func checkoutCopy(targetRepo vcs.RepoWriter, copyBranchName, sourceWorkspaceID, preCommitID string) (string, error) {
	err := targetRepo.FetchBranch(copyBranchName, sourceWorkspaceID)
	if err != nil {
		return "", fmt.Errorf("failed to fetch on target: %w", err)
	}

	err = targetRepo.CreateBranchTrackingUpstream(copyBranchName)
	if err != nil {
		return "", fmt.Errorf("failed to create branch on target: %w", err)
	}

	err = targetRepo.CheckoutBranchWithForce(copyBranchName)
	if err != nil {
		return "", fmt.Errorf("failed to checkout branch on target: %w", err)
	}

	err = targetRepo.CreateBranchTrackingUpstream(sourceWorkspaceID)
	if err != nil {
		return "", fmt.Errorf("failed to create workspace branch on target: %w", err)
	}

	// reset to the parent commit (to remove the commit, but keep the changes)
	err = targetRepo.ResetMixed(preCommitID)
	if err != nil {
		return "", fmt.Errorf("failed to restore to parent on target: %w", err)
	}

	// checkout the workspace branch,
	err = targetRepo.CheckoutBranchSafely(sourceWorkspaceID)
	if err != nil {
		return "", fmt.Errorf("failed to checkout workspace on target: %w", err)
	}

	return copyBranchName, nil
}
