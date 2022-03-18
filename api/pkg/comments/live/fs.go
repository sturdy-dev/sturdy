package live

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"

	"getsturdy.com/api/pkg/changes"
	"getsturdy.com/api/pkg/codebases"
	"getsturdy.com/api/pkg/snapshots"
	db_snapshots "getsturdy.com/api/pkg/snapshots/db"
	"getsturdy.com/api/pkg/workspaces"
	"getsturdy.com/api/vcs"
	"getsturdy.com/api/vcs/executor"
)

var ErrNoFiles = fmt.Errorf("files not found")

func WorkspaceFS(
	executorProvider executor.Provider,
	snapshotsRepo db_snapshots.Repository,
	ws *workspaces.Workspace,
	// newLines directly links to comment.LineIsNew
	newLines bool,
) (fs.FS, error) {
	if ws.ViewID != nil {
		return viewFS(executorProvider, ws.CodebaseID, *ws.ViewID, newLines), nil
	}
	if ws.LatestSnapshotID != nil {
		snapshot, err := snapshotsRepo.Get(*ws.LatestSnapshotID)
		if err != nil {
			return nil, fmt.Errorf("failed to get snapshot: %w", err)
		}
		return snapshotFS(executorProvider, snapshot, newLines), nil
	}
	return nil, ErrNoFiles
}

func ChangeFS(executorProvider executor.Provider, change *changes.Change, newLines bool) (fs.FS, error) {
	return &changeFS{
		executorProvider: executorProvider,
		change:           change,
		newLines:         newLines,
	}, nil
}

type changeFS struct {
	executorProvider executor.Provider
	change           *changes.Change
	// newLines directly links to comment.LineIsNew
	newLines bool
}

func (cfs *changeFS) Open(path string) (fs.File, error) {
	if cfs.change.CommitID == nil {
		return nil, fmt.Errorf("commit id is nil")
	}
	var file fs.File
	return file, cfs.executorProvider.New().GitRead(func(repo vcs.RepoGitReader) error {
		if cfs.newLines {
			var err error
			file, err = fileFromCommit(repo, *cfs.change.CommitID, path)
			return err
		}

		parents, err := repo.GetCommitParents(*cfs.change.CommitID)
		if err != nil {
			return fmt.Errorf("failed to get commit parents: %w", err)
		}
		if len(parents) != 1 {
			return fmt.Errorf("expected 1 parent, got %d", len(parents))
		}
		file, err = fileFromCommit(repo, parents[0], path)
		return err
	}).ExecTrunk(cfs.change.CodebaseID, fmt.Sprintf("open %s", path))
}

type SnapshotsFS struct {
	executorProvider executor.Provider
	snapshot         *snapshots.Snapshot
	// newLines directly links to comment.LineIsNew
	newLines bool
}

func snapshotFS(
	executorProvider executor.Provider,
	snapshot *snapshots.Snapshot,
	newLines bool,
) *SnapshotsFS {
	return &SnapshotsFS{
		executorProvider: executorProvider,
		snapshot:         snapshot,
		newLines:         newLines,
	}
}

func (snapshotsFS *SnapshotsFS) Open(path string) (fs.File, error) {
	var file fs.File
	return file, snapshotsFS.executorProvider.New().GitRead(func(repo vcs.RepoGitReader) error {
		if snapshotsFS.newLines {
			var err error
			file, err = fileFromCommit(repo, snapshotsFS.snapshot.CommitID, path)
			return err
		}

		parents, err := repo.GetCommitParents(snapshotsFS.snapshot.CommitID)
		if err != nil {
			return fmt.Errorf("failed to get commit parents: %w", err)
		}
		if len(parents) != 1 {
			return fmt.Errorf("expected 1 parent, got %d", len(parents))
		}
		file, err = fileFromCommit(repo, parents[0], path)
		return err
	}).ExecTrunk(snapshotsFS.snapshot.CodebaseID, fmt.Sprintf("open %s", path))
}

type ViewFS struct {
	executorProvider executor.Provider
	viewID           string
	codebaseID       codebases.ID
	// newLines directly links to comment.LineIsNew
	newLines bool
}

func viewFS(
	executorProvider executor.Provider,
	codebaseID codebases.ID,
	viewID string,
	newLines bool,
) fs.FS {
	return &ViewFS{
		executorProvider: executorProvider,
		viewID:           viewID,
		codebaseID:       codebaseID,
		newLines:         newLines,
	}
}

func (viewFS *ViewFS) Open(path string) (fs.File, error) {
	if viewFS.newLines {
		var file fs.File
		return file, viewFS.executorProvider.New().Read(func(repo vcs.RepoReader) error {
			osFile, err := os.Open(filepath.Join(repo.Path(), path))
			if err != nil {
				return fmt.Errorf("failed to open file: %w", err)
			}
			file = osFile
			return nil
		}).ExecView(viewFS.codebaseID, viewFS.viewID, "viewFsFromDisk")
	}

	var file fs.File
	return file, viewFS.executorProvider.New().GitRead(func(repo vcs.RepoGitReader) error {
		headCommit, err := repo.HeadCommit()
		if err != nil {
			return fmt.Errorf("failed to get head commit: %w", err)
		}
		defer headCommit.Free()
		file, err = fileFromCommit(repo, headCommit.Id().String(), path)
		return err
	}).ExecView(viewFS.codebaseID, viewFS.viewID, "viewFsFromCommit")
}

func fileFromCommit(repo vcs.RepoGitReader, commitSHA string, path string) (fs.File, error) {
	blob, err := repo.FileBlobAtCommit(commitSHA, path)
	switch {
	case err == nil:
		defer blob.Free()
		return fileFrom(bytes.NewReader(blob.Contents())), nil
	case errors.Is(err, vcs.ErrFileNotFound):
		return nil, fs.ErrNotExist
	default:
		return nil, fmt.Errorf("failed to get file contents: %w", err)
	}
}

type memoryFile struct {
	readCloser io.ReadCloser
}

func fileFrom(reader io.Reader) fs.File {
	return &memoryFile{
		readCloser: io.NopCloser(reader),
	}
}

func (f *memoryFile) Stat() (fs.FileInfo, error) {
	return nil, fmt.Errorf("not implemented")
}

func (f *memoryFile) Read(b []byte) (int, error) {
	return f.readCloser.Read(b)
}

func (f *memoryFile) Close() error {
	return f.readCloser.Close()
}
