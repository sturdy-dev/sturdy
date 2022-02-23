package live

import (
	"database/sql"
	"errors"
	"fmt"
	"io/fs"
	"io/ioutil"
	"strings"
	"unicode"

	"getsturdy.com/api/pkg/comments"
	db_comments "getsturdy.com/api/pkg/comments/db"
	db_snapshots "getsturdy.com/api/pkg/snapshots/db"
	"getsturdy.com/api/pkg/workspaces"
	"getsturdy.com/api/vcs/executor"
)

func contextMatchCount(rows []string, startRowNum int, context []string, preprocess func(string) string) (fuzzyMatches int) {
	// for each row in the context, see if any row matches fuzzily
contextLoop:
	for d := 0; d < len(context) && startRowNum+d < len(rows); d++ {
		for fuzzyDelta := 0; fuzzyDelta < len(context) && startRowNum+fuzzyDelta < len(rows); fuzzyDelta++ {
			if preprocess(rows[startRowNum+fuzzyDelta]) == preprocess(context[d]) {
				fuzzyMatches++
				continue contextLoop
			}
		}
	}
	return
}

func fuzzyNewLocation(c comments.Comment, fileContents string) int {
	rows := strings.Split(fileContents, "\n")
	context := strings.Split(*c.Context, "\n")

	// Max 5 lines of context
	if len(context) == 6 {
		context = context[:5]
	}

	contextStartsAt := *c.ContextStartsAtLine
	contextDelta := c.LineStart - contextStartsAt

	var bestFuzzyMatches = 0
	var bestFuzzyMatchesDelta = 0
	var remainingIters = -1

	// Search with whitespace normalization
	for delta := 0; delta < 1000; delta++ {

		// Once we find a partial fuzzy match, keep searching for 10 more iterations, in hope to find an even better match
		if bestFuzzyMatches > 0 {
			if remainingIters == 0 {
				break
			} else if remainingIters > 0 {
				remainingIters--
			}
		}

		if contextStartsAt-delta >= 0 {
			fuzzyMatches := contextMatchCount(rows, contextStartsAt-delta, context, normalizeWhitespace)
			if fuzzyMatches > bestFuzzyMatches {
				if remainingIters == -1 && fuzzyMatches > 2 { // We found a first fuzzy match, search up to 10 more times
					remainingIters = 10
				}
				bestFuzzyMatches = fuzzyMatches
				bestFuzzyMatchesDelta = -delta
			}
		}

		if contextStartsAt+delta < len(rows) {
			fuzzyMatches := contextMatchCount(rows, contextStartsAt+delta, context, normalizeWhitespace)
			if fuzzyMatches > bestFuzzyMatches {
				if remainingIters == -1 && fuzzyMatches > 2 { // We found a first fuzzy match, search up to 10 more times
					remainingIters = 10
				}
				bestFuzzyMatches = fuzzyMatches
				bestFuzzyMatchesDelta = delta
			}
		}

		// We have a perfect match, stop searching
		if bestFuzzyMatches == 5 {
			break
		}
	}

	// Threshold is min(3, len(rows))
	threshold := 3
	if threshold > len(rows) {
		threshold = len(rows)
	}

	// at least 3 out of 5 lines matched
	if bestFuzzyMatches >= threshold {
		return contextStartsAt + bestFuzzyMatchesDelta + contextDelta + 1
	}

	return -1
}

func noop(str string) string {
	return str
}

func normalizeWhitespace(str string) string {
	var b strings.Builder
	b.Grow(len(str))
	for _, ch := range str {
		if !unicode.IsSpace(ch) {
			b.WriteRune(ch)
		}
	}
	return b.String()
}

func GetWorkspaceComments(
	commentRepo db_comments.Repository,
	ws *workspaces.Workspace,
	executorProvider executor.Provider,
	snapshotRepo db_snapshots.Repository,
) ([]comments.Comment, error) {
	comms, err := commentRepo.GetByWorkspace(ws.ID)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("could not get comments by workspace: %w", err)
	}

	newFilesFS, err := WorkspaceFS(executorProvider, snapshotRepo, ws, true)
	switch {
	case err == nil:
	case errors.Is(err, ErrNoFiles):
		return nil, nil
	default:
		return nil, fmt.Errorf("could not prepare workspace filesystem: %w", err)
	}

	oldFilesFS, err := WorkspaceFS(executorProvider, snapshotRepo, ws, false)
	if err != nil {
		return nil, fmt.Errorf("could not prepare workspace filesystem: %w", err)
	}

	// fuzzily update line numbers
	for i, c := range comms {
		if c.Context == nil || c.ContextStartsAtLine == nil {
			comms[i].LineStart = -1
			comms[i].LineEnd = -1
			continue
		}

		var (
			file fs.File
			err  error
		)
		if c.LineIsNew {
			file, err = newFilesFS.Open(c.Path)
		} else {
			if c.OldPath != nil {
				file, err = oldFilesFS.Open(*c.OldPath)
			} else {
				file, err = oldFilesFS.Open(c.Path)
			}
		}
		switch {
		case err == nil:
		case errors.Is(err, fs.ErrNotExist), errors.Is(err, executor.ErrIsRebasing):
			comms[i].LineStart = -1
			comms[i].LineEnd = -1
			continue
		default:
			return nil, fmt.Errorf("could not open file %s: %w", c.Path, err)
		}

		contents, err := ioutil.ReadAll(file)
		if err != nil {
			return nil, fmt.Errorf("could not read file contents: %w", err)
		}

		newLoc := fuzzyNewLocation(c, string(contents))
		comms[i].LineStart = newLoc
		comms[i].LineEnd = newLoc

		if err := file.Close(); err != nil {
			return nil, fmt.Errorf("could not close file: %w", err)
		}
	}

	return comms, nil
}
