package vcs

import (
	"time"

	git "github.com/libgit2/git2go/v33"
)

type LogEntry struct {
	Time             time.Time
	RawCommitMessage string
	ID               string // This is a Commit ID

	Name  string
	Email string

	// If this commits has been landed on sturdytrunk (potentially with another commit ID)
	// TODO(gustav): delete this field
	IsLanded bool
}

func (repo *repository) LogHead(limit int) ([]*LogEntry, error) {
	defer getMeterFunc("LogHead")()
	revwalk, err := repo.r.Walk()
	if err != nil {
		return nil, err
	}
	defer revwalk.Free()

	err = revwalk.PushHead()
	if err != nil {
		return nil, err
	}

	return repo.log(revwalk, limit)
}

func (repo *repository) LogBranch(branchName string, limit int) ([]*LogEntry, error) {
	defer getMeterFunc("LogBranch")()
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

	err = revwalk.Push(branch.Target())
	if err != nil {
		return nil, err
	}

	return repo.log(revwalk, limit)
}

func CommitLogEntry(commit *git.Commit) *LogEntry {
	id := commit.Id().String()
	committer := commit.Committer()

	return &LogEntry{
		Time:             committer.When,
		RawCommitMessage: commit.Message(),
		ID:               id,
		Name:             committer.Name,
		Email:            committer.Email,
	}
}

func (repo *repository) log(revwalk *git.RevWalk, limit int) ([]*LogEntry, error) {
	defer getMeterFunc("log")()
	var out []*LogEntry

	var i int
	revwalk.Iterate(func(commit *git.Commit) bool {
		// If merge commit, skip it from the list
		if commit.ParentCount() >= 2 {
			return true
		}
		out = append(out, CommitLogEntry(commit))

		// Break at limit
		i++
		if i == limit {
			return false
		}

		return true
	})

	return out, nil
}
