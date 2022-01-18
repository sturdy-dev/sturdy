package vcs

import (
	"bytes"
	"fmt"
	"os/exec"

	git "github.com/libgit2/git2go/v33"
)

// ApplyPatchesToIndex patches parsed and applied sequentially and the index is written afterwards.
// If there is an error with any of the patches, the function returns before writing the index.
// Returns a treeID or an error. The treeID can be passed to CommitIndexTree.
func (repo *repository) ApplyPatchesToIndex(patches [][]byte) (*git.Oid, error) {
	defer getMeterFunc("ApplyPatchesToIndex")()

	opts, err := git.DefaultApplyOptions()
	if err != nil {
		return nil, err
	}

	for _, p := range patches {
		diff, err := git.DiffFromBuffer(p, repo.r)
		if err != nil {
			return nil, fmt.Errorf("failed to parse change: %w", err)
		}

		err = repo.r.ApplyDiff(diff, git.ApplyLocationIndex, opts)
		if err != nil {
			_ = diff.Free()
			return nil, fmt.Errorf("failed to apply change: %w", err)
		}
		_ = diff.Free()
	}
	index, err := repo.r.Index()
	if err != nil {
		return nil, fmt.Errorf("failed to access vcs index: %w", err)
	}
	defer index.Free()

	oid, err := index.WriteTree()
	if err != nil {
		return nil, fmt.Errorf("write tree failed: %w", err)
	}
	return oid, nil
}

func (repo *repository) ApplyPatchesToWorkdir(patches [][]byte) error {
	defer getMeterFunc("ApplyPatchesToWorkdir")()
	for _, patch := range patches {
		cmd := exec.Command("git", "apply", "-")
		cmd.Dir = repo.path
		cmd.Stdin = bytes.NewReader(patch)

		out, err := cmd.CombinedOutput()
		if err != nil {
			return fmt.Errorf("unexpected error: '%s' %w", string(out), err)
		}
	}
	return nil
}

func (r *repository) CanApplyPatch(patch []byte) (bool, error) {
	defer getMeterFunc("CanApplyPatch")()
	cmd := exec.Command("git", "apply", "--check", "-")
	cmd.Dir = r.path
	cmd.Stdin = bytes.NewReader(patch)

	out, err := cmd.CombinedOutput()
	if err != nil {
		if bytes.HasPrefix(out, []byte("error: patch failed: ")) {
			return false, nil
		}
		if bytes.HasSuffix(out, []byte(": already exists in working directory\n")) {
			return false, nil
		}
		if bytes.HasSuffix(out, []byte(": No such file or directory\n")) {
			return false, nil
		}
		return false, fmt.Errorf("unexpected error: '%s' %w", string(out), err)
	}
	return true, nil
}
