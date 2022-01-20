package vcs

import (
	git "github.com/libgit2/git2go/v33"
)

func (r *repository) ResetMixed(commitID string) error {
	defer getMeterFunc("ResetMixed")()
	oid, err := git.NewOid(commitID)
	if err != nil {
		return err
	}

	commit, err := r.r.LookupCommit(oid)
	if err != nil {
		return err
	}
	defer commit.Free()

	err = r.r.ResetToCommit(commit, git.ResetMixed, &git.CheckoutOptions{
		Strategy: git.CheckoutSafe,
	})
	if err != nil {
		return err
	}

	return nil
}

func (r *repository) ResetHard(commitID string) error {
	defer getMeterFunc("ResetHard")()
	oid, err := git.NewOid(commitID)
	if err != nil {
		return err
	}

	commit, err := r.r.LookupCommit(oid)
	if err != nil {
		return err
	}
	defer commit.Free()

	err = r.r.ResetToCommit(commit, git.ResetHard, &git.CheckoutOptions{
		Strategy: git.CheckoutForce,
	})
	if err != nil {
		return err
	}

	return nil
}
