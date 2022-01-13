package vcs

import (
	"fmt"
	"os/exec"

	git "github.com/libgit2/git2go/v33"
)

func (r *repository) FetchOriginCLI() error {
	defer getMeterFunc("FetchOriginCLI")()
	cmd := exec.Command("git", "fetch", "origin")
	cmd.Dir = r.path
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed '%s': %w", string(output), err)
	}
	return nil
}

// InitRebase starts a rebasing operation
// Returns an error in case of unexpected failures, conflicts that require user intervention
// are not unexpected, and are indicated via SturdyRebase.Action, which will be RebaseHaveConflicts.
func (r *repository) InitRebase(ontoRemoteName, ontoBranchName string) (*SturdyRebase, []RebasedCommit, error) {
	defer getMeterFunc("InitRebase")()
	headBranch, err := r.HeadBranch()
	if err != nil {
		return nil, nil, fmt.Errorf("failed to find rebase head: %w", err)
	}
	headBranchCommit, err := r.annotatedCommitFromBranchName(headBranch)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to find rebase head annotated commit: %w", err)
	}
	defer headBranchCommit.Free()
	onto, err := r.RemoteBranchCommit(ontoRemoteName, ontoBranchName)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to find target: %w", err)
	}
	defer onto.Free()
	annotatedOntoCommit, err := r.r.LookupAnnotatedCommit(onto.Id())
	if err != nil {
		return nil, nil, fmt.Errorf("failed to find target: %w", err)
	}
	defer annotatedOntoCommit.Free()
	// Find common ancestor (this is the rebase "upstream")
	commonAncestor, err := r.r.MergeBase(
		headBranchCommit.Id(),
		onto.Id(),
	)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to find common ancestor: %w", err)
	}
	annotatedCommonAncestor, err := r.r.LookupAnnotatedCommit(commonAncestor)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to find common ancestor: %w", err)
	}
	defer annotatedCommonAncestor.Free()
	// Stash unsaved changes before attempting rebase
	err = r.stashUnsavedForRebase()
	if err != nil {
		return nil, nil, fmt.Errorf("stashing failed: %w", err)
	}
	// Start rebase
	rebase, err := r.r.InitRebase(
		headBranchCommit,
		annotatedCommonAncestor,
		annotatedOntoCommit,
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

func (r *repository) GetCommit(id string) (*git.Commit, error) {
	defer getMeterFunc("GetCommit")()
	oid, err := git.NewOid(id)
	if err != nil {
		return nil, err
	}
	commit, err := r.r.LookupCommit(oid)
	if err != nil {
		return nil, err
	}
	return commit, nil
}

func (repo *repository) LogBranchUntilTrunk(branchName string, limit int) ([]*LogEntry, error) {
	defer getMeterFunc("LogBranchUntilTrunk")()
	branch, err := repo.r.LookupBranch(branchName, git.BranchLocal)
	if err != nil {
		return nil, err
	}
	defer branch.Free()
	revwalk, err := repo.r.Walk()
	if err != nil {
		return nil, err
	}
	defer revwalk.Free()
	err = revwalk.PushRange("refs/heads/sturdytrunk.." + branch.Reference.Name())
	if err != nil {
		return nil, err
	}
	return repo.log(revwalk, limit)
}
