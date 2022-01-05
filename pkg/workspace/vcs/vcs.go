package vcs

import (
	"fmt"

	"mash/vcs"
	"mash/vcs/provider"
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

func ListChanges(trunkProvider provider.TrunkProvider, codebaseID, workspaceID string, limit int) ([]*vcs.LogEntry, error) {
	r, err := trunkProvider.TrunkRepo(codebaseID)
	if err != nil {
		return nil, fmt.Errorf("failed find codebaseID %s: %w", codebaseID, err)
	}
	defer r.Free()

	log, err := r.LogBranch(workspaceID, limit)
	if err != nil {
		return nil, fmt.Errorf("failed fetch changes for workspaceID %s: %w", workspaceID, err)
	}

	revlist, err := r.RevlistCherryPickLeftOnly(workspaceID, "sturdytrunk")
	if err != nil {
		return nil, fmt.Errorf("failed to get revlist: %w", err)
	}
	unmergedCommits := make(map[string]struct{})
	for _, c := range revlist {
		unmergedCommits[c] = struct{}{}
	}

	for k, v := range log {
		if _, ok := unmergedCommits[v.ID]; !ok {
			v.IsLanded = true
			log[k] = v
		}
	}

	// Remove the root commit from this log
	// This is stupid
	// TODO: Can we remove the "Root commit" stuff?
	if log[len(log)-1].RawCommitMessage == "Root Commit" {
		log = log[0 : len(log)-1]
	}

	return log, nil
}
