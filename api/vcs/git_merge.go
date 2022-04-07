package vcs

import (
	"fmt"
	"time"

	git "github.com/libgit2/git2go/v33"
)

// MergeBranches returns a new index with the merge of the two branches
// The returned index must be freed by the caller.
func (r *repository) MergeBranches(ourBranchName, theirBranchName string) (*git.Index, error) {
	defer getMeterFunc("MergeBranches")()

	ourBranch, err := r.r.LookupBranch(ourBranchName, git.BranchAll)
	if err != nil {
		return nil, fmt.Errorf("failed to look up our branch %s: %w", ourBranchName, err)
	}

	ourCommit, err := r.r.LookupCommit(ourBranch.Target())
	if err != nil {
		return nil, fmt.Errorf("failed to look up our commit: %w", err)
	}

	theirBranch, err := r.r.LookupBranch(theirBranchName, git.BranchAll)
	if err != nil {
		return nil, fmt.Errorf("failed to look up their branch %s: %w", theirBranchName, err)
	}

	theirCommit, err := r.r.LookupCommit(theirBranch.Target())
	if err != nil {
		return nil, fmt.Errorf("failed to look up their commit: %w", err)
	}

	opts, err := git.DefaultMergeOptions()
	if err != nil {
		return nil, err
	}

	idx, err := r.r.MergeCommits(ourCommit, theirCommit, &opts)
	if err != nil {
		return nil, err
	}

	return idx, nil
}

func (r *repository) MergeBranchInto(branchName, mergeIntoBranchName string) (mergeCommitID string, err error) {
	defer getMeterFunc("MergeBranchInto")()
	sourceBranch, err := r.r.LookupBranch(branchName, git.BranchLocal)
	if err != nil {
		return "", fmt.Errorf("failed to look up branch %s: %w", branchName, err)
	}

	destinationBranch, err := r.r.LookupBranch(mergeIntoBranchName, git.BranchLocal)
	if err != nil {
		return "", fmt.Errorf("failed to look up branch %s: %w", mergeIntoBranchName, err)
	}

	sourceCommit, err := r.r.LookupCommit(sourceBranch.Target())
	if err != nil {
		return "", fmt.Errorf("failed to get sourceCommit: %w", err)
	}

	destinationCommit, err := r.r.LookupCommit(destinationBranch.Target())
	if err != nil {
		return "", fmt.Errorf("failed to get destinationCommit: %w", err)
	}

	sourceTree, err := sourceCommit.Tree()
	if err != nil {
		return "", fmt.Errorf("failed to get sourceTree: %w", err)
	}

	destinationTree, err := sourceCommit.Tree()
	if err != nil {
		return "", fmt.Errorf("failed to get destinationTree: %w", err)
	}

	opts, err := git.DefaultMergeOptions()
	if err != nil {
		return "", err
	}
	opts.TreeFlags = git.MergeTreeFailOnConflict

	mergeBase, err := r.r.MergeBase(sourceCommit.Id(), destinationCommit.Id())
	if err != nil {
		return "", fmt.Errorf("failed to get mergeBase: %w", err)
	}

	mergeBaseCommit, err := r.r.LookupCommit(mergeBase)
	if err != nil {
		return "", fmt.Errorf("failed to get mergeBaseCommit: %w", err)
	}

	mergeBaseTree, err := mergeBaseCommit.Tree()
	if err != nil {
		return "", fmt.Errorf("failed to get mergeBaseTree: %w", err)
	}

	idx, err := r.r.MergeTrees(mergeBaseTree, sourceTree, destinationTree, &opts)
	if err != nil {
		return "", fmt.Errorf("failed to mergeTrees: %w", err)
	}

	treeOid, err := idx.WriteTreeTo(r.r)
	if err != nil {
		return "", fmt.Errorf("failed to write new tree: %w", err)
	}

	mergedTree, err := r.r.LookupTree(treeOid)
	if err != nil {
		return "", fmt.Errorf("failed to get mergedTree: %w", err)
	}

	sig := git.Signature{Name: "merge", Email: "merge@getsturdy.com", When: time.Now()}
	mergeCommit, err := r.r.CreateCommit("refs/heads/"+mergeIntoBranchName, &sig, &sig,
		fmt.Sprintf("Merge %s into %s", mergeIntoBranchName, branchName),
		mergedTree,
		// Parents
		destinationCommit, sourceCommit)
	if err != nil {
		return "", err
	}

	return mergeCommit.String(), nil
}

func (r *repository) CommonAncestor(commitA, commitB string) (string, error) {
	idA, err := git.NewOid(commitA)
	if err != nil {
		return "", fmt.Errorf("could not parse commitA: %w", err)
	}
	idB, err := git.NewOid(commitB)
	if err != nil {
		return "", fmt.Errorf("could not parse commitB: %w", err)
	}

	mergeBase, err := r.r.MergeBase(idA, idB)
	if err != nil {
		return "", fmt.Errorf("failed to get mergeBase: %w", err)
	}

	return mergeBase.String(), nil
}
