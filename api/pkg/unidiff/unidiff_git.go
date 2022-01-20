package unidiff

import (
	"fmt"
	"io"

	git "github.com/libgit2/git2go/v33"
)

type gitPatchReader struct {
	numDeltas int
	idx       int
	diff      *git.Diff
	err       error
}

func NewGitPatchReader(diff *git.Diff) PatchReader {
	nd, err := diff.NumDeltas()
	if err != nil {
		// Store the error, and return it on the first read
		return &gitPatchReader{err: err}
	}
	return &gitPatchReader{
		idx:       0,
		numDeltas: nd,
		diff:      diff,
	}
}

func (g *gitPatchReader) ReadPatch() (string, error) {
	if g.err != nil {
		return "", g.err
	}

	if g.idx == g.numDeltas {
		return "", io.EOF
	}
	patch, err := g.diff.Patch(g.idx)
	if err != nil {
		return "", fmt.Errorf("failed to get diff patch: %w", err)
	}

	text, err := patch.String()
	if err != nil {
		return "", fmt.Errorf("failed to get patch as string: %w", err)
	}

	err = patch.Free()
	if err != nil {
		return "", fmt.Errorf("failed to free patch: %w", err)
	}

	g.idx++

	return text, nil
}
