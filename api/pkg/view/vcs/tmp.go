package vcs

import (
	"fmt"
	"os"

	"mash/vcs"
	"mash/vcs/provider"

	"github.com/google/uuid"
)

// TemporaryViewFromSnapshotWithID return a temporary view with a snapshot on top of it.
// Caller must call the cancel function.
func TemporaryViewFromSnapshotWithID(
	repoProvider provider.RepoProvider,
	viewID string,
	codebaseID string,
	workspaceID string,
	snapshotID string,
) (vcs.RepoWriter, func() error, error) {
	if err := Create(repoProvider, codebaseID, workspaceID, viewID); err != nil {
		return nil, nil, fmt.Errorf("failed to create tmp view: %w", err)
	}

	viewRepo, err := repoProvider.ViewRepo(codebaseID, viewID)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to open tmp view: %w", err)
	}

	// Copy snapshot onto the view
	snapshotBranchName := fmt.Sprintf("snapshot-%s", snapshotID)
	if _, err := CheckoutWorkspaceSnapshot(repoProvider, codebaseID, workspaceID, snapshotBranchName, viewID); err != nil {
		return nil, nil, fmt.Errorf("failed to checkout snapshot: %w", err)
	}

	return viewRepo, func() error {
		return os.RemoveAll(repoProvider.ViewPath(codebaseID, viewID))
	}, nil

}

// TemporaryViewFromSnapshot is the same as TemporaryViewFromSnapshotWithID, but it generates the view id.
func TemporaryViewFromSnapshot(repoProvider provider.RepoProvider, codebaseID, workspaceID, snapshotID string) (vcs.RepoWriter, func() error, error) {
	viewID := fmt.Sprintf("tmp-%s", uuid.New().String())
	return TemporaryViewFromSnapshotWithID(repoProvider, viewID, codebaseID, workspaceID, snapshotID)
}

func TemporaryView(
	repoProvider provider.RepoProvider,
	codebaseID string,
	checkoutBranchName string,
) (vcs.RepoWriter, func() error, error) {
	return TemporaryViewWithID(repoProvider, fmt.Sprintf("tmp-%s", uuid.New().String()), codebaseID, checkoutBranchName)
}

func TemporaryViewWithID(
	repoProvider provider.RepoProvider,
	viewID string,
	codebaseID string,
	checkoutBranchName string,
) (vcs.RepoWriter, func() error, error) {
	if err := Create(repoProvider, codebaseID, checkoutBranchName, viewID); err != nil {
		return nil, nil, fmt.Errorf("failed to create tmp view: %w", err)
	}

	viewRepo, err := repoProvider.ViewRepo(codebaseID, viewID)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to open tmp view: %w", err)
	}

	return viewRepo, func() error {
		return os.RemoveAll(repoProvider.ViewPath(codebaseID, viewID))
	}, nil

}
