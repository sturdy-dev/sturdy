package vcs

import (
	"fmt"

	git "github.com/libgit2/git2go/v33"
)

func (r *repository) GetCommitParents(commitID string) ([]string, error) {
	defer getMeterFunc("GetCommitParents")()
	oid, err := git.NewOid(commitID)
	if err != nil {
		return nil, err
	}

	commit, err := r.r.LookupCommit(oid)
	if err != nil {
		return nil, err
	}
	defer commit.Free()

	var parents []string
	pc := commit.ParentCount()
	for i := uint(0); i < pc; i++ {
		p := commit.Parent(i)
		parents = append(parents, p.Id().String())
		p.Free()
	}

	return parents, nil
}

func (r *repository) CommitMessage(id string) (author *git.Signature, message string, err error) {
	defer getMeterFunc("CommitMessage")()
	co, err := r.r.RevparseSingle(id)
	if err != nil {
		return nil, "", err
	}
	defer co.Free()

	c, err := co.AsCommit()
	if err != nil {
		return nil, "", err
	}
	defer c.Free()

	return c.Author(), c.Message(), nil
}

func (r *repository) ShowCommit(id string) (diffs []string, entry *LogEntry, err error) {
	defer getMeterFunc("ShowCommit")()
	co, err := r.r.RevparseSingle(id)
	if err != nil {
		return nil, nil, err
	}
	defer co.Free()

	c, err := co.AsCommit()
	if err != nil {
		return nil, nil, err
	}
	defer c.Free()

	cTree, err := c.Tree()
	if err != nil {
		return nil, nil, err
	}
	defer cTree.Free()

	// Diff against parent if a parent exists
	p := c.Parent(0)
	var pTree *git.Tree
	if p != nil {
		pTree, err = p.Tree()
		if err != nil {
			return nil, nil, err
		}
		defer pTree.Free()
	}

	diff, err := r.r.DiffTreeToTree(pTree, cTree, nil)
	if err != nil {
		return nil, nil, err
	}
	defer diff.Free()

	err = sturdyFindSimilar(diff)
	if err != nil {
		return nil, nil, err
	}

	numDeltas, err := diff.NumDeltas()
	if err != nil {
		return nil, nil, err
	}

	var out []string
	for i := 0; i < numDeltas; i++ {
		patch, _ := diff.Patch(i)
		text, _ := patch.String()
		out = append(out, text)
		_ = patch.Free()
	}

	return out, CommitLogEntry(c), nil
}

func (r *repository) BranchHasCommit(branchName, commitID string) (bool, error) {
	defer getMeterFunc("BranchHasCommit")()
	branch, err := r.r.LookupBranch(branchName, git.BranchLocal)
	if err != nil {
		return false, err
	}
	defer branch.Free()

	branchTarget := branch.Target()

	commitOid, err := git.NewOid(commitID)
	if err != nil {
		return false, err
	}

	// The branch HEAD is this commit
	if branchTarget.Equal(commitOid) {
		return true, nil
	}

	descendant, err := r.r.DescendantOf(branchTarget, commitOid)
	if err != nil {
		return false, err
	}

	return descendant, nil
}

type FileContents struct {
	Path     string
	Contents []byte
}

func (r *repository) CreateCommitWithFiles(files []FileContents, newBranchName string) (string, error) {
	headTree, err := r.getHead()
	if err != nil {
		return "", fmt.Errorf("failed to get head tree: %w", err)
	}

	headCommit, err := r.HeadCommit()
	if err != nil {
		return "", fmt.Errorf("failed to get head commit: %w", err)
	}

	tb, err := r.r.TreeBuilderFromTree(headTree)
	if err != nil {
		return "", fmt.Errorf("failed to get tree builder: %w", err)
	}

	for _, fc := range files {
		blobOid, err := r.r.CreateBlobFromBuffer(fc.Contents)
		if err != nil {
			return "", fmt.Errorf("failed to create blob: %w", err)
		}

		if err := tb.Insert(fc.Path, blobOid, git.FilemodeBlob); err != nil {
			return "", fmt.Errorf("failed to insert blob: %w", err)
		}
	}

	newTreeOid, err := tb.Write()
	if err != nil {
		return "", fmt.Errorf("failed to write tree: %w", err)
	}

	newTree, err := r.r.LookupTree(newTreeOid)
	if err != nil {
		return "", fmt.Errorf("failed to lookup new tree: %w", err)
	}

	signature := git.Signature{Email: "noreply@getsturdy.com", Name: "Sturdy"}

	newCommit, err := r.r.CreateCommit("refs/heads/"+newBranchName, &signature, &signature, "CreateCommitWithFiles", newTree, headCommit)
	if err != nil {
		return "", fmt.Errorf("failed to commit: %w", err)
	}

	return newCommit.String(), nil
}
