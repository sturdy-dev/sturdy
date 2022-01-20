package vcs

import (
	"fmt"

	git "github.com/libgit2/git2go/v33"
)

func (r *repository) RevertOnBranch(revertCommitID, branchName string) (string, error) {
	defer getMeterFunc("RevertOnBranch")()
	id, err := git.NewOid(revertCommitID)
	if err != nil {
		return "", err
	}
	revertCommit, err := r.r.LookupCommit(id)
	if err != nil {
		return "", err
	}
	defer revertCommit.Free()

	branchHeadCommitId, err := r.BranchCommitID(branchName)
	if err != nil {
		return "", err
	}
	branchHeadCommitOid, err := git.NewOid(branchHeadCommitId)
	if err != nil {
		return "", err
	}
	branchHeadCommit, err := r.r.LookupCommit(branchHeadCommitOid)
	if err != nil {
		return "", err
	}

	opts, err := git.DefaultMergeOptions()
	if err != nil {
		return "", err
	}

	idx, err := r.r.RevertCommit(revertCommit, branchHeadCommit, 0, &opts)
	if err != nil {
		return "", err
	}
	defer idx.Free()

	treeID, err := idx.WriteTreeTo(r.r)
	if err != nil {
		return "", err
	}

	tree, err := r.r.LookupTree(treeID)
	if err != nil {
		return "", fmt.Errorf("lookup tree failed: %w", err)
	}
	defer tree.Free()

	branch, err := r.r.Head()
	if err != nil {
		return "", fmt.Errorf("get head failed: %w", err)
	}
	defer branch.Free()

	signature := git.Signature{Email: "noreply@getsturdy.com", Name: "Sturdy"}

	oid, err := r.r.CreateCommit("refs/heads/"+branchName, &signature, &signature, "revert", tree, branchHeadCommit)
	if err != nil {
		return "", fmt.Errorf("failed to create commit: %w", err)
	}
	return oid.String(), nil
}
