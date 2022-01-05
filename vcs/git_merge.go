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
	ourRef, err := r.r.References.Lookup("refs/remotes/origin/" + ourBranchName)
	if err != nil {
		return nil, fmt.Errorf("failed to look up reference %s: %w", ourBranchName, err)
	}
	defer ourRef.Free()

	ourCommit, err := r.r.LookupCommit(ourRef.Branch().Target())
	if err != nil {
		return nil, fmt.Errorf("failed to look up commit: %w", err)
	}
	defer ourCommit.Free()

	theirRef, err := r.r.References.Lookup("refs/remotes/origin/" + theirBranchName)
	if err != nil {
		return nil, fmt.Errorf("failed to look up reference %s: %w", theirBranchName, err)
	}
	defer theirRef.Free()

	theirCommit, err := r.r.LookupCommit(theirRef.Branch().Target())
	if err != nil {
		return nil, fmt.Errorf("failed to look up commit: %w", err)
	}
	defer theirCommit.Free()

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

func (r *repository) MergeBranchInto(branchName, mergeIntoBranchName string) error {
	defer getMeterFunc("MergeBranchInto")()
	sourceBranch, err := r.r.LookupBranch(branchName, git.BranchLocal)
	if err != nil {
		return fmt.Errorf("failed to look up branch %s: %w", branchName, err)
	}
	defer sourceBranch.Free()

	destinationBranch, err := r.r.LookupBranch(mergeIntoBranchName, git.BranchLocal)
	if err != nil {
		return fmt.Errorf("failed to look up branch %s: %w", mergeIntoBranchName, err)
	}
	defer destinationBranch.Free()

	sourceAnnotatedCommit, err := r.r.AnnotatedCommitFromRef(sourceBranch.Reference)
	if err != nil {
		return fmt.Errorf("failed to get sourceAnnotatedCommit: %w", err)
	}
	defer sourceAnnotatedCommit.Free()

	destinationAnnotatedCommit, err := r.r.AnnotatedCommitFromRef(destinationBranch.Reference)
	if err != nil {
		return fmt.Errorf("failed to get destinationAnnotatedCommit: %w", err)
	}
	defer destinationAnnotatedCommit.Free()

	sourceCommit, err := r.r.LookupCommit(sourceBranch.Target())
	if err != nil {
		return fmt.Errorf("failed to get sourceCommit: %w", err)
	}
	defer sourceCommit.Free()

	destinationCommit, err := r.r.LookupCommit(destinationBranch.Target())
	if err != nil {
		return fmt.Errorf("failed to get destinationCommit: %w", err)
	}
	defer destinationCommit.Free()

	sourceTree, err := sourceCommit.Tree()
	if err != nil {
		return fmt.Errorf("failed to get sourceTree: %w", err)
	}
	defer sourceTree.Free()

	destinationTree, err := sourceCommit.Tree()
	if err != nil {
		return fmt.Errorf("failed to get destinationTree: %w", err)
	}
	defer destinationTree.Free()

	opts, err := git.DefaultMergeOptions()
	if err != nil {
		return err
	}
	opts.TreeFlags = git.MergeTreeFailOnConflict

	mergeBase, err := r.r.MergeBase(sourceCommit.Id(), destinationCommit.Id())
	if err != nil {
		return fmt.Errorf("failed to get mergeBase: %w", err)
	}

	mergeBaseCommit, err := r.r.LookupCommit(mergeBase)
	if err != nil {
		return fmt.Errorf("failed to get mergeBaseCommit: %w", err)
	}
	defer mergeBaseCommit.Free()

	mergeBaseTree, err := mergeBaseCommit.Tree()
	if err != nil {
		return fmt.Errorf("failed to get mergeBaseTree: %w", err)
	}
	defer mergeBaseTree.Free()

	idx, err := r.r.MergeTrees(mergeBaseTree, sourceTree, destinationTree, &opts)
	if err != nil {
		return fmt.Errorf("failed to mergeTrees: %w", err)
	}
	defer idx.Free()

	treeOid, err := idx.WriteTreeTo(r.r)
	if err != nil {
		return fmt.Errorf("failed to write new tree: %w", err)
	}

	mergedTree, err := r.r.LookupTree(treeOid)
	if err != nil {
		return fmt.Errorf("failed to get mergedTree: %w", err)
	}
	defer mergedTree.Free()

	sig := git.Signature{Name: "merge", Email: "merge@getsturdy.com", When: time.Now()}
	_, err = r.r.CreateCommit("refs/heads/"+mergeIntoBranchName, &sig, &sig,
		fmt.Sprintf("Merge %s into %s", mergeIntoBranchName, branchName),
		mergedTree,
		// Parents
		destinationCommit, sourceCommit)
	if err != nil {
		return err
	}

	return nil
}
