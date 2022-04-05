package service

import (
	"context"
	"fmt"
	"io/fs"
	"io/ioutil"
	"path"

	"github.com/h2non/filetype"

	"getsturdy.com/api/pkg/comments/live"
	"getsturdy.com/api/pkg/file"
	db_snapshots "getsturdy.com/api/pkg/snapshots/db"
	"getsturdy.com/api/pkg/workspaces"
	"getsturdy.com/api/vcs/executor"
)

type Service struct {
	executorProvider executor.Provider
	snapshotsRepo    db_snapshots.Repository
}

func New(executorProvider executor.Provider) *Service {
	return &Service{
		executorProvider: executorProvider,
	}
}

func (s *Service) ReadWorkspaceFile(ctx context.Context, ws *workspaces.Workspace, filePath string, isNew bool) ([]byte, error) {
	if filePath != path.Clean(filePath) {
		return nil, fmt.Errorf("unexpected path: %s", filePath)
	}

	fs, err := live.WorkspaceFS(s.executorProvider, s.snapshotsRepo, ws, isNew)
	if err != nil {
		return nil, fmt.Errorf("failed to create fs: %w", err)
	}

	fp, err := fs.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}

	data, err := ioutil.ReadAll(fp)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	return data, nil
}

var fileExtFilter = map[string]struct{}{
	"jpg":  {},
	"jpeg": {},
	"png":  {},
	"gif":  {},
	"webp": {},
}

func (s *Service) WorkspaceFileType(ctx context.Context, ws *workspaces.Workspace, filePath string, isNew bool) (file.Type, error) {
	// ext is not in filter, don't validate it
	ext := path.Ext(filePath)
	if _, ok := fileExtFilter[ext]; !ok {
		return file.UnknownType, nil
	}

	fs, err := live.WorkspaceFS(s.executorProvider, s.snapshotsRepo, ws, isNew)
	if err != nil {
		return file.UnknownType, fmt.Errorf("failed to create fs: %w", err)
	}

	fp, err := fs.Open(filePath)
	if err != nil {
		return file.UnknownType, fmt.Errorf("failed to open file: %w", err)
	}

	fileType, err := s.detectFileType(fp)
	if err != nil {
		return file.UnknownType, fmt.Errorf("failed to detect file type: %w", err)
	}

	return fileType, nil
}

func (s *Service) detectFileType(fp fs.File) (file.Type, error) {
	// We only have to pass the file header = first 261 bytes
	head := make([]byte, 261)
	if _, err := fp.Read(head); err != nil {
		return file.UnknownType, fmt.Errorf("failed to read file header: %w", err)
	}

	if filetype.IsImage(head) {
		return file.ImageType, nil
	}

	// TODO: detect text?

	return file.BinaryType, nil
}
