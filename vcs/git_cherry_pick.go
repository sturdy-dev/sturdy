package vcs

import (
	"errors"
	"fmt"
	"os/exec"
	"strings"

	git "github.com/libgit2/git2go/v33"
)

var cherryPickOptions = git.CherrypickOptions{
	MergeOptions: git.MergeOptions{
		// Defaults from https://github.com/libgit2/libgit2sharp/blob/master/LibGit2Sharp/MergeOptionsBase.cs
		RenameThreshold: 50,
		TargetLimit:     200,
	},
	CheckoutOptions: git.CheckoutOptions{
		Strategy: git.CheckoutConflictStyleMerge, // CheckoutSafe
	},
}

func (r *repository) CherryPickOnto(commitID, onto string) (newCommitID string, conflicted bool, conflictingFiles []string, err error) {
	defer getMeterFunc("CherryPickOnto")()
	ontoID, err := git.NewOid(onto)
	if err != nil {
		return "", false, nil, fmt.Errorf("failed to parse commit ID: %w", err)
	}
	pickOntoCommit, err := r.r.LookupCommit(ontoID)
	if err != nil {
		return "", false, nil, fmt.Errorf("failed to find commit: %w", err)
	}
	defer pickOntoCommit.Free()
	id, err := git.NewOid(commitID)
	if err != nil {
		return "", false, nil, fmt.Errorf("failed to parse commit ID: %w", err)
	}
	commit, err := r.r.LookupCommit(id)
	if err != nil {
		return "", false, nil, fmt.Errorf("failed to find commit: %w", err)
	}
	defer commit.Free()
	newIdx, err := r.r.CherrypickCommit(commit, pickOntoCommit, cherryPickOptions)
	if err != nil {
		return "", false, nil, fmt.Errorf("cherry picking failed: %w", err)
	}
	defer newIdx.Free()

	if newIdx.HasConflicts() {
		conflictingFiles, err := ConflictingFilesInIndex(newIdx)
		if err != nil {
			return "", false, nil, fmt.Errorf("failed to get conflicting files: %w", err)
		}
		return "", true, conflictingFiles, nil
	}

	treeOid, err := newIdx.WriteTreeTo(r.r)
	if err != nil {
		return "", false, nil, fmt.Errorf("write tree failed: %w", err)
	}

	author := commit.Author()
	newCommitID, err = r.CommitIndexTree(treeOid, commit.Message(), *author)
	if err != nil {
		return "", false, nil, fmt.Errorf("CommitIndexTree failed: %w", err)
	}

	return newCommitID, false, nil, nil
}

// base = the commit (and parents that you're rebasing)
// When you're rebasing, --left-only will return the commits that exist on base that don't exist on onto
func (r *repository) RevlistCherryPickLeftOnly(base, onto string) ([]string, error) {
	defer getMeterFunc("RevlistCherryPickLeftOnly")()
	return r.revlistCherryPick(base, onto, []string{"--left-only"})
}

// When you're rebasing, --left-only will return the commits that exist on onto, that don't exist on base
func (r *repository) RevlistCherryPickRightOnly(base, onto string) ([]string, error) {
	defer getMeterFunc("RevlistCherryPickRightOnly")()
	return r.revlistCherryPick(base, onto, []string{"--right-only"})
}

func (r *repository) revlistCherryPick(base, onto string, extraArgs []string) ([]string, error) {
	defer getMeterFunc("revlistCherryPick")()
	args := []string{"rev-list", "--cherry-pick"}
	args = append(args, extraArgs...)
	args = append(args, fmt.Sprintf("%s...%s", base, onto))

	checkoutCmd := exec.Command("git", args...)
	checkoutCmd.Dir = r.path

	output, err := checkoutCmd.Output()
	if err != nil {
		return nil, err
	}

	var commits []string
	for _, c := range strings.Split(strings.TrimSpace(string(output)), "\n") {
		if len(c) == 40 {
			commits = append(commits, c)
		}
	}

	// reverse the list, so that it's returned in the same order as the commits would be picked when being applied
	for i, j := 0, len(commits)-1; i < j; i, j = i+1, j-1 {
		commits[i], commits[j] = commits[j], commits[i]
	}

	return commits, nil
}

func (r *repository) BranchCommitID(branchName string) (string, error) {
	defer getMeterFunc("BranchCommitID")()
	branch, err := r.r.LookupBranch(branchName, git.BranchLocal)
	if err != nil {
		return "", fmt.Errorf("failed to look up branch: %w", err)
	}
	defer branch.Free()
	commit, err := r.r.LookupCommit(branch.Target())
	if err != nil {
		return "", fmt.Errorf("failed to look up commit: %w", err)
	}
	defer commit.Free()
	return commit.Id().String(), nil
}

var ErrCommitNotFound = errors.New("commit not found")

func (r *repository) BranchFirstNonMergeCommit(branchName string) (string, error) {
	defer getMeterFunc("BranchFirstNonMergeCommit")()
	branch, err := r.r.LookupBranch(branchName, git.BranchLocal)
	if err != nil {
		return "", fmt.Errorf("failed to look up branch: %w", err)
	}
	defer branch.Free()

	// Breadth first search
	q := []*git.Oid{branch.Target()}
	for {
		if len(q) == 0 {
			return "", ErrCommitNotFound
		}

		id := q[0]
		q = q[1:]

		commit, err := r.r.LookupCommit(id)
		if err != nil {
			return "", fmt.Errorf("failed to look up commit: %w", err)
		}

		// Non-merge commit
		if commit.ParentCount() == 1 {
			return commit.Id().String(), nil
		}

		// Add to queue
		for parentId := uint(0); parentId < commit.ParentCount(); parentId++ {
			q = append(q, commit.ParentId(parentId))
		}
	}
}

func (r *repository) InitRebaseRaw(head, onto string) (*SturdyRebase, []RebasedCommit, error) {
	defer getMeterFunc("InitRebaseRaw")()
	// Stash unsaved changes before attempting rebase
	err := r.stashUnsavedForRebase()
	if err != nil {
		return nil, nil, fmt.Errorf("stashing failed: %w", err)
	}

	headID, err := git.NewOid(head)
	if err != nil {
		return nil, nil, err
	}
	headAnnotated, err := r.r.LookupAnnotatedCommit(headID)
	if err != nil {
		return nil, nil, err
	}
	defer headAnnotated.Free()

	headCommit, err := r.r.LookupCommit(headID)
	if err != nil {
		return nil, nil, err
	}
	defer headCommit.Free()

	headParent := headCommit.Parent(0)
	defer headParent.Free()

	headParentAnnoteted, err := r.r.LookupAnnotatedCommit(headParent.Id())
	if err != nil {
		return nil, nil, err
	}
	defer headParentAnnoteted.Free()

	ontoID, err := git.NewOid(onto)
	if err != nil {
		return nil, nil, err
	}
	ontoAnnotated, err := r.r.LookupAnnotatedCommit(ontoID)
	if err != nil {
		return nil, nil, err
	}
	defer ontoAnnotated.Free()

	// Start rebase
	rebase, err := r.r.InitRebase(
		headAnnotated,
		headParentAnnoteted,
		ontoAnnotated,
		commonRebaseOptions,
	)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to init rebase: %w", err)
	}

	status := &SturdyRebase{
		repo:      r,
		gitRebase: rebase,
	}

	_, rebasedCommits, err := status.Continue()
	if err != nil {
		return nil, nil, fmt.Errorf("rebasing resolving failed: %w", err)
	}

	return status, rebasedCommits, nil
}
