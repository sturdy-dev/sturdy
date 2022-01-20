package vcs

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strings"
	"time"

	git "github.com/libgit2/git2go/v33"
)

type SturdyRebaseStatus int

const (
	RebaseInProgress    SturdyRebaseStatus = iota // Used during initialization
	RebaseCanContinue                             // The rebase can continue
	RebaseHaveConflicts                           // There are conflicts that must be resolved before continuing
	RebaseCompleted                               // The rebase is completed
)

var commonRebaseOptions = &git.RebaseOptions{
	InMemory: 0,
	MergeOptions: git.MergeOptions{
		// Defaults from https://github.com/libgit2/libgit2sharp/blob/master/LibGit2Sharp/MergeOptionsBase.cs
		RenameThreshold: 50,
		TargetLimit:     200,
	},
	CheckoutOptions: git.CheckoutOptions{
		Strategy: git.CheckoutSafe,
	},
}

type SturdyRebase struct {
	repo      *repository
	gitRebase *git.Rebase

	completed           bool
	lastCompletedCommit string
}

var NoRebaseInProgress = errors.New("no rebasing in progress")

// OpenRebase resumes a rebasing operation started by InitRebase
func (r *repository) OpenRebase() (*SturdyRebase, error) {
	defer getMeterFunc("OpenRebase")()
	reb, err := r.r.OpenRebase(commonRebaseOptions)
	if err != nil {
		if gitErr, ok := err.(*git.GitError); ok && gitErr.Class == git.ErrClassRebase && gitErr.Code == git.ErrNotFound {
			return nil, NoRebaseInProgress
		}
		return nil, fmt.Errorf("could not open rebase: %w", err)
	}
	return &SturdyRebase{
		repo:      r,
		gitRebase: reb,
	}, nil
}

type SturdyRebaseResolve struct {
	Path    string
	Version string
}

func (rebase *SturdyRebase) ResolveFiles(resolves []SturdyRebaseResolve) error {
	idx, err := rebase.repo.r.Index()
	if err != nil {
		return fmt.Errorf("failed to get index in rebase: %w", err)
	}
	defer idx.Free()

	for _, resolve := range resolves {
		err := rebase.resolveFile(idx, resolve.Path, resolve.Version)
		if err != nil {
			return err
		}
	}

	err = idx.Write()
	if err != nil {
		return fmt.Errorf("failed to write new index: %w", err)
	}

	return nil
}

func (rebase *SturdyRebase) resolveFile(idx *git.Index, filePath string, version string) error {
	conflict, err := idx.Conflict(filePath)
	if err != nil {
		return fmt.Errorf("failed to get conflict: %w", err)
	}

	// The perspective when rebasing is from the new bases pov
	var use *git.IndexEntry
	switch version {
	case "custom":
		// Add the file as is
		err = idx.RemoveConflict(filePath)
		if err != nil {
			return fmt.Errorf("failed to remove conflict from index: %w", err)
		}
		err = idx.AddByPath(filePath)
		if err != nil {
			return fmt.Errorf("failed to add resolved file: %w", err)
		}
		return nil
	case "workspace":
		use = conflict.Their
	case "trunk":
		use = conflict.Our
	default:
		return fmt.Errorf("unknown version: %s", version)
	}

	fullFilePath := path.Join(rebase.repo.path, filePath)

	// use is nil when the file that we're picking is deleted
	if use != nil {
		blb, err := rebase.repo.r.LookupBlob(use.Id)
		if err != nil {
			return fmt.Errorf("failed to get blob: %w", err)
		}
		defer blb.Free()

		err = ioutil.WriteFile(fullFilePath, blb.Contents(), os.FileMode(use.Mode))
		if err != nil {
			return fmt.Errorf("failed to write resolution: %w", err)
		}
	} else {
		// Make sure to delete the file
		if err := os.Remove(fullFilePath); err != nil {
			return fmt.Errorf("failed to remove file during conflict resolution: %w", err)
		}
	}

	if err := idx.RemoveConflict(filePath); err != nil {
		return fmt.Errorf("failed to remove conflict from index: %w", err)
	}

	if use != nil {
		if err := idx.AddByPath(filePath); err != nil {
			return fmt.Errorf("failed to add resolved file: %w", err)
		}
	} else {
		if err := idx.RemoveByPath(filePath); err != nil {
			return fmt.Errorf("failed to remove resolved file: %w", err)
		}
	}

	return nil
}

func (rebase *SturdyRebase) ConflictingFiles() ([]string, error) {
	idx, err := rebase.repo.r.Index()
	if err != nil {
		return nil, fmt.Errorf("failed to get index in rebase: %w", err)
	}
	defer idx.Free()
	return ConflictingFilesInIndex(idx)
}

func ConflictingFilesInIndex(idx *git.Index) ([]string, error) {
	conflicts, err := idx.ConflictIterator()
	if err != nil {
		return nil, fmt.Errorf("failed to parse conflicts: %w", err)
	}

	var paths []string

	for {
		c, err := conflicts.Next()
		if err != nil {
			if gitErr, ok := err.(*git.GitError); ok && gitErr.Code == git.ErrIterOver {
				break
			}
			return nil, fmt.Errorf("failed to get next conflict: %w", err)
		}
		if c.Their != nil {
			paths = append(paths, c.Their.Path)
		} else if c.Our != nil {
			paths = append(paths, c.Our.Path)
		} else if c.Ancestor != nil {
			paths = append(paths, c.Ancestor.Path)
		}
	}

	return paths, nil
}

type RebasedCommit struct {
	OldCommitID string
	NewCommitID string
	Noop        bool
}

// Continue the rebase until the next conflict, or the rebase is completed
func (rebase *SturdyRebase) Continue() (conflicts bool, rebasedCommits []RebasedCommit, err error) {
	idx, err := rebase.repo.r.Index()
	if err != nil {
		return false, nil, fmt.Errorf("failed to get index: %w", err)
	}
	defer idx.Free()

	for {
		var operation *git.RebaseOperation

		rebaseOpIdx, err := rebase.gitRebase.CurrentOperationIndex()
		if err != nil {
			if errors.Is(err, git.ErrRebaseNoOperation) {
				// Continue with the next operation
				operation, err = rebase.gitRebase.Next()
				if err != nil {
					if gitErr, ok := err.(*git.GitError); ok && gitErr.Code == git.ErrIterOver {
						break
					}
					return false, nil, fmt.Errorf("failed to run next operation: %w", err)
				}
			} else {
				return false, nil, fmt.Errorf("failed to get operation: %w", err)
			}
		} else {
			operation = rebase.gitRebase.OperationAt(rebaseOpIdx)
		}

		if idx.HasConflicts() {
			return true, rebasedCommits, nil
		}

		// Load commit data so that we can re-use it
		commit, err := rebase.repo.r.LookupCommit(operation.Id)
		if err != nil {
			return false, nil, fmt.Errorf("failed to find commit in rebase: %w", err)
		}
		defer commit.Free()

		// Add all unstaged changes to the index (to track new updates, deletes, and added files _during_ the sync)
		if _, err := rebase.repo.AddFilesToIndex([]string{"."}); err != nil {
			return false, nil, fmt.Errorf("failed to add unstaged to index: %w", err)
		}

		var noop bool
		err = rebase.gitRebase.Commit(operation.Id, commit.Author(), commit.Committer(), commit.Message())
		if err != nil {
			if gitErr, ok := err.(*git.GitError); ok && gitErr.Class == git.ErrClassRebase && gitErr.Code == git.ErrApplied {
				noop = true
			} else {
				return false, nil, fmt.Errorf("failed to commit: %w", err)
			}
		}

		if noop {
			rebasedCommits = append(rebasedCommits, RebasedCommit{
				OldCommitID: commit.Id().String(),
				Noop:        noop,
			})
		} else {
			headCommitAfter, err := rebase.repo.HeadCommit()
			if err != nil {
				return false, nil, fmt.Errorf("failed to get HEAD after commit: %w", err)
			}
			// Record which commit resulted in what new commit
			rebasedCommits = append(rebasedCommits, RebasedCommit{
				OldCommitID: commit.Id().String(),
				NewCommitID: headCommitAfter.Id().String(),
			})
			headCommitAfter.Free()
		}

		rebase.lastCompletedCommit = commit.Id().String()

		// Advance to the next operation
		_, err = rebase.gitRebase.Next()
		if err != nil {
			if gitErr, ok := err.(*git.GitError); ok && gitErr.Code == git.ErrIterOver {
				break
			}
			return false, nil, fmt.Errorf("failed to run next operation: %w", err)
		}
	}

	err = rebase.gitRebase.Finish()
	if err != nil {
		return false, nil, fmt.Errorf("failed to complete rebase: %w", err)
	}

	// Stash unsaved changes before attempting rebase
	err = rebase.repo.stashPopFromRebase()
	if err != nil {
		return false, nil, fmt.Errorf("failed to pop stash: %w", err)
	}

	rebase.completed = true

	return false, rebasedCommits, nil
}

func (rebase *SturdyRebase) LastCompletedCommit() string {
	return rebase.lastCompletedCommit
}

func (rebase *SturdyRebase) Status() (SturdyRebaseStatus, error) {
	if rebase.completed {
		return RebaseCompleted, nil
	}
	// This is not a nice solution
	// TODO: Create an interface with "Status() (SturdyRebaseStatus, error)" so that
	//       this case can be handled nicely when SturdyRebase is "empty" (happens in a rebase with dropped commits)
	if rebase.repo == nil {
		return RebaseCompleted, nil
	}

	idx, err := rebase.repo.r.Index()
	if err != nil {
		return 0, fmt.Errorf("failed to get index: %w", err)
	}
	defer idx.Free()

	if idx.HasConflicts() {
		return RebaseHaveConflicts, nil
	}

	return RebaseCanContinue, nil
}

func (rebase *SturdyRebase) Progress() (current, total uint, err error) {
	current, err = rebase.gitRebase.CurrentOperationIndex()
	if err != nil {
		return 0, 0, fmt.Errorf("failed to get current operation index: %w", err)
	}
	total = rebase.gitRebase.OperationCount()
	return
}

type ConflictDiffs struct {
	WorkspacePatch string
	TrunkPatch     string
}

func (rebase *SturdyRebase) ConflictDiff(filePath string) (ConflictDiffs, error) {
	idx, err := rebase.repo.r.Index()
	if err != nil {
		return ConflictDiffs{}, fmt.Errorf("failed to get index: %w", err)
	}
	defer idx.Free()

	conflict, err := idx.Conflict(filePath)
	if err != nil {
		return ConflictDiffs{}, fmt.Errorf("failed to get conflict: %w", err)
	}

	var diffs ConflictDiffs

	diffs.WorkspacePatch, err = rebase.patchBetweenIndexEntries(conflict.Ancestor, conflict.Their)
	if err != nil {
		return ConflictDiffs{}, fmt.Errorf("failed to create workspace patch: %w", err)
	}

	diffs.TrunkPatch, err = rebase.patchBetweenIndexEntries(conflict.Ancestor, conflict.Our)
	if err != nil {
		return ConflictDiffs{}, fmt.Errorf("failed to create trunk patch: %w", err)
	}

	return diffs, nil
}

func (rebase *SturdyRebase) patchBetweenIndexEntries(entryA, entryB *git.IndexEntry) (string, error) {
	var blobA, blobB *git.Blob
	var blobApath, blobBpath string
	var err error

	// It's ok for one of entryA or entryB to be nil

	if entryA != nil {
		blobApath = entryA.Path
		blobA, err = rebase.repo.r.LookupBlob(entryA.Id)
		if err != nil {
			return "", fmt.Errorf("failed to get blob: %w", err)
		}
		defer blobA.Free()
	}

	if entryB != nil {
		blobBpath = entryB.Path
		blobB, err = rebase.repo.r.LookupBlob(entryB.Id)
		if err != nil {
			return "", fmt.Errorf("failed to get blob: %w", err)
		}
		defer blobB.Free()
	}

	var patch string

	err = git.DiffBlobs(
		blobA,
		blobApath,
		blobB,
		blobBpath,
		nil,
		func(delta git.DiffDelta, progress float64) (git.DiffForEachHunkCallback, error) {
			// An example diff:
			// diff --git b/abc.txt a/abc.txt
			// index 0000000..a3bb749 100644
			// --- b/abc.txt
			// +++ a/abc.txt
			// @@ -1,6 +1,7 @@
			//  a
			//  b
			//  d
			// +e
			//  f
			//  g
			//  h

			nameOrDevNull := func(prefix, name string) string {
				if len(name) == 0 {
					return "/dev/null"
				}
				return prefix + name
			}

			oldName := nameOrDevNull("a/", delta.OldFile.Path)
			newName := nameOrDevNull("b/", delta.NewFile.Path)

			patch = fmt.Sprintf("diff --git %s %s\n", oldName, newName)
			patch += fmt.Sprintf("index 0000000..0000000 %x\n", delta.NewFile.Mode)
			patch += fmt.Sprintf("--- %s\n", oldName)
			patch += fmt.Sprintf("+++ %s\n", newName)

			return func(hunk git.DiffHunk) (git.DiffForEachLineCallback, error) {
				patch += fmt.Sprintf("%s", hunk.Header)

				return func(line git.DiffLine) error {
					var prefix rune

					switch line.Origin {
					case git.DiffLineAddition, git.DiffLineAddEOFNL:
						prefix = '+'
					case git.DiffLineDeletion, git.DiffLineDelEOFNL:
						prefix = '-'
					case git.DiffLineContext, git.DiffLineContextEOFNL:
						prefix = ' '
					default:
						return fmt.Errorf("unknown line.Origin: %+v", line)
					}

					patch += fmt.Sprintf("%c%s", prefix, line.Content)
					return nil
				}, nil
			}, nil
		},
		git.DiffDetailLines,
	)
	if err != nil {
		return "", fmt.Errorf("patch creation failed: %w", err)
	}

	return patch, nil
}

func (r *repository) stashUnsavedForRebase() error {
	defer getMeterFunc("stashUnsavedForRebase")()
	_, err := r.r.Stashes.Save(&git.Signature{
		Name:  "Stasher",
		Email: "stash@getsturdy.com",
		When:  time.Now(),
	}, "Stashing to perform rebase", git.StashDefault)
	if err != nil {
		// Nothing to stash
		if gitErr := err.(*git.GitError); gitErr.Class == git.ErrClassStash && gitErr.Code == git.ErrNotFound {
			return nil
		}
		return err
	}
	return nil
}

var errStopLoop = errors.New("stop loop")

func (r *repository) stashPopFromRebase() error {
	defer getMeterFunc("stashPopFromRebase")()
	err := r.r.Stashes.Foreach(func(index int, message string, id *git.Oid) error {
		if strings.Contains(message, "Stashing to perform rebase") {
			err := r.r.Stashes.Apply(index, git.StashApplyOptions{
				CheckoutOptions: git.CheckoutOpts{
					Strategy: git.CheckoutSafe,
				},
			})
			if err != nil {
				return fmt.Errorf("failed to apply stash: %w", err)
			}
			return errStopLoop
		}
		return nil
	})
	if err != nil && !errors.Is(err, errStopLoop) {
		return err
	}
	return nil
}
