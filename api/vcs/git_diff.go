package vcs

import (
	"fmt"

	"getsturdy.com/api/vcs/diff"

	git "github.com/libgit2/git2go/v33"
)

type DiffOptions struct {
	gitMaxSize  int
	withIndex   bool
	withReverse bool
}

type DiffOption func(*DiffOptions)

func WithGitMaxSize(size int) DiffOption {
	return func(opts *DiffOptions) {
		opts.gitMaxSize = size
	}
}

func WithIndex() DiffOption {
	return func(opts *DiffOptions) {
		opts.withIndex = true
	}
}

func WithReverse() DiffOption {
	return func(opts *DiffOptions) {
		opts.withReverse = true
	}
}

func getDiffOptions(opts ...DiffOption) *DiffOptions {
	options := &DiffOptions{
		gitMaxSize: 50_000_000, // 50MB by default
	}
	for _, applyOption := range opts {
		applyOption(options)
	}
	return options
}

func (r *repository) CurrentDiff(opts ...DiffOption) (*git.Diff, error) {
	defer getMeterFunc("CurrentDiff")()
	headTree, err := r.getHead()
	if err != nil {
		return nil, fmt.Errorf("CurrentDiff could not get HEAD: %w", err)
	}
	defer headTree.Free()

	opts = append(opts, WithIndex())

	return r.workdirDiffAgainstTree(headTree, opts...)
}

func (r *repository) CurrentDiffNoIndex() (*git.Diff, error) {
	defer getMeterFunc("CurrentDiffNoIndex")()

	headTree, err := r.getHead()
	if err != nil {
		return nil, fmt.Errorf("CurrentDiffNoIndex could not get HEAD: %w", err)
	}
	defer headTree.Free()

	return r.workdirDiffAgainstTree(headTree)
}

func (r *repository) workdirDiffAgainstTree(diffAgainst *git.Tree, opts ...DiffOption) (*git.Diff, error) {
	defer getMeterFunc("workdirDiffAgainstTree")()

	o := getDiffOptions(opts...)

	gitOpts, err := git.DefaultDiffOptions()
	if err != nil {
		return nil, err
	}

	gitOpts.Flags = git.DiffNormal | git.DiffIncludeUntracked | git.DiffRecurseUntracked | git.DiffShowUntrackedContent
	// opts.Flags |= git.DiffIgnoreWhitespace | git.DiffIgnoreWhitespaceChange | git.DiffIgnoreWitespaceEol
	if o.withReverse {
		gitOpts.Flags |= git.DiffReverse
	}

	gitOpts.OldPrefix = diff.DiffOldPrefix
	gitOpts.NewPrefix = diff.DiffNewPrefix

	gitOpts.MaxSize = o.gitMaxSize

	var diff *git.Diff
	if o.withIndex {
		diff, err = r.r.DiffTreeToWorkdirWithIndex(diffAgainst, &gitOpts)
	} else {
		diff, err = r.r.DiffTreeToWorkdir(diffAgainst, &gitOpts)
	}
	if err != nil {
		return nil, fmt.Errorf("failed to perform diffing: %w", err)
	}

	err = sturdyFindSimilar(diff)
	if err != nil {
		return nil, err
	}

	return diff, nil
}

func (repo *repository) DiffCommits(firstCommitID, secondCommitID string) (*git.Diff, error) {
	defer getMeterFunc("DiffCommits")()
	firstCommitOID, err := git.NewOid(firstCommitID)
	if err != nil {
		return nil, err
	}
	secondCommitOID, err := git.NewOid(secondCommitID)
	if err != nil {
		return nil, err
	}

	firstCommit, err := repo.r.LookupCommit(firstCommitOID)
	if err != nil {
		return nil, err
	}
	defer firstCommit.Free()
	secondCommit, err := repo.r.LookupCommit(secondCommitOID)
	if err != nil {
		return nil, err
	}
	defer secondCommit.Free()

	firstCommitTree, err := firstCommit.Tree()
	if err != nil {
		return nil, err
	}
	defer firstCommitTree.Free()
	secondCommitTree, err := secondCommit.Tree()
	if err != nil {
		return nil, err
	}
	defer secondCommitTree.Free()

	diff, err := repo.r.DiffTreeToTree(firstCommitTree, secondCommitTree, nil)
	if err != nil {
		return nil, err
	}

	err = sturdyFindSimilar(diff)
	if err != nil {
		return nil, err
	}

	return diff, nil
}

func (repo *repository) DiffCommitToRoot(firstCommitID string) (*git.Diff, error) {
	defer getMeterFunc("DiffCommitToRoot")()
	firstCommitOID, err := git.NewOid(firstCommitID)
	if err != nil {
		return nil, err
	}

	firstCommit, err := repo.r.LookupCommit(firstCommitOID)
	if err != nil {
		return nil, err
	}

	firstCommitTree, err := firstCommit.Tree()
	if err != nil {
		return nil, err
	}

	diff, err := repo.r.DiffTreeToTree(nil, firstCommitTree, nil)
	if err != nil {
		return nil, err
	}

	err = sturdyFindSimilar(diff)
	if err != nil {
		return nil, err
	}

	return diff, nil
}

func sturdyFindSimilar(diff *git.Diff) error {
	opts, err := git.DefaultDiffFindOptions()
	if err != nil {
		return fmt.Errorf("could not get default find opts: %w", err)
	}
	opts.Flags |= git.DiffFindRenames |
		git.DiffFindIgnoreWhitespace |
		git.DiffFindForUntracked |
		git.DiffFindRemoveUnmodified // Remove _unmodified_ hunks from the result

	// All available flags:
	//
	// DiffFindByConfig
	// DiffFindRenames
	// DiffFindRenamesFromRewrites
	// DiffFindCopies
	// DiffFindCopiesFromUnmodified
	// DiffFindRewrites
	// DiffFindBreakRewrites
	// DiffFindAndBreakRewrites
	// DiffFindForUntracked
	// DiffFindAll
	// DiffFindIgnoreLeadingWhitespace
	// DiffFindIgnoreWhitespace
	// DiffFindDontIgnoreWhitespace
	// DiffFindExactMatchOnly
	// DiffFindBreakRewritesForRenamesOnly
	// DiffFindRemoveUnmodified

	opts.RenameThreshold = 50
	opts.CopyThreshold = 50

	err = diff.FindSimilar(&opts)
	if err != nil {
		return fmt.Errorf("find similar within diff failed: %w", err)
	}
	return nil
}
