package vcs

import (
	"bufio"
	"io/fs"
	"strings"

	"getsturdy.com/api/pkg/changes"
	"getsturdy.com/api/pkg/comments/live"
	db_snapshots "getsturdy.com/api/pkg/snapshots/db"
	"getsturdy.com/api/pkg/workspaces"
	"getsturdy.com/api/vcs/executor"
)

func GetChangeContext(lineNumber int, lineIsNew bool, filePath string, oldFilePath *string, change *changes.Change, executorProvider executor.Provider) (context string, startsAtLine int, err error) {
	changeFs, err := live.ChangeFS(executorProvider, change, lineIsNew)
	if err != nil {
		return "", -1, err
	}
	return getFSContext(lineNumber, lineIsNew, filePath, oldFilePath, changeFs)
}

func GetWorkspaceContext(lineNumber int, lineIsNew bool, filePath string, oldFilePath *string, ws *workspaces.Workspace, executorProvider executor.Provider, snapshotRepo db_snapshots.Repository) (context string, startsAtLine int, err error) {
	workspaceFS, err := live.WorkspaceFS(executorProvider, snapshotRepo, ws, lineIsNew)
	if err != nil {
		return "", -1, err
	}
	return getFSContext(lineNumber, lineIsNew, filePath, oldFilePath, workspaceFS)
}

func getFSContext(lineNumber int, lineIsNew bool, filePath string, oldFilePath *string, filesystem fs.FS) (context string, startsAtLine int, err error) {
	var file fs.File
	if lineIsNew {
		// New lines, always use the new file name
		file, err = filesystem.Open(filePath)
	} else {
		if oldFilePath != nil {
			// Comment on old line, and oldFilePath is set
			file, err = filesystem.Open(*oldFilePath)
		} else {
			// Comment on old line, and oldFilePath is not set
			file, err = filesystem.Open(filePath)
		}
	}
	if err != nil {
		return "", -1, err
	}
	defer file.Close()

	s := bufio.NewScanner(file)
	contextLines := 2 // two lines before and after
	var res strings.Builder

	l := 1
	startsAtLine = -1

	for ; s.Scan(); l++ {
		if l >= lineNumber-contextLines && l <= lineNumber+contextLines {
			if startsAtLine == -1 {
				startsAtLine = l
			}
			res.WriteString(s.Text() + "\n")
		}
		if l == lineNumber+contextLines {
			break
		}
	}
	return res.String(), startsAtLine, nil
}
