package service

import (
	"bytes"
	"context"
	"crypto/sha1"
	"fmt"
	"io"
	"io/fs"
	"io/ioutil"
	"net/http"
	"path"
	"strings"

	"github.com/h2non/filetype"

	"getsturdy.com/api/pkg/changes"
	"getsturdy.com/api/pkg/codebases"
	"getsturdy.com/api/pkg/comments/live"
	"getsturdy.com/api/pkg/file"
	db_snapshots "getsturdy.com/api/pkg/snapshots/db"
	"getsturdy.com/api/pkg/workspaces"
	"getsturdy.com/api/vcs/executor"
	provider "getsturdy.com/api/vcs/provider/configuration"
)

type Service struct {
	executorProvider executor.Provider
	snapshotsRepo    db_snapshots.Repository
	vcsConfiguration *provider.Configuration
}

func New(
	executorProvider executor.Provider,
	snapshotsRepo db_snapshots.Repository,
	vcsConfiguration *provider.Configuration,
) *Service {
	return &Service{
		executorProvider: executorProvider,
		snapshotsRepo:    snapshotsRepo,
		vcsConfiguration: vcsConfiguration,
	}
}

func (s *Service) ReadWorkspaceFile(ctx context.Context, ws *workspaces.Workspace, filePath string, isNew bool) ([]byte, error) {
	fsys, err := live.WorkspaceFS(s.executorProvider, s.snapshotsRepo, ws, isNew)
	if err != nil {
		return nil, fmt.Errorf("failed to create fs: %w", err)
	}

	return s.readFile(fsys, filePath, ws.CodebaseID)
}

func (s *Service) ReadChangeFile(ctx context.Context, ch *changes.Change, filePath string, isNew bool) ([]byte, error) {
	fsys, err := live.ChangeFS(s.executorProvider, ch, isNew)
	if err != nil {
		return nil, fmt.Errorf("failed to create fs: %w", err)
	}

	return s.readFile(fsys, filePath, ch.CodebaseID)
}

func (s *Service) readFile(fsys fs.FS, filePath string, codebaseID codebases.ID) ([]byte, error) {
	if filePath != path.Clean(filePath) {
		return nil, fmt.Errorf("unexpected path: %s", filePath)
	}

	fp, err := fsys.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}

	data, err := ioutil.ReadAll(fp)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	// if the file LFS pointer, serve the file from the LFS server
	if len(data) < 2000 && bytes.HasPrefix(data, []byte("version https://git-lfs.github.com/spec/v1\n")) {
		lines := strings.Split(string(data), "\n")
		if len(lines) < 2 || len(lines[1]) != 75 {
			return nil, fmt.Errorf("failed to parse lfs file")
		}
		sha := lines[1][len("oid sha256:"):]
		resp, err := http.Get("http://" + s.vcsConfiguration.LFS.Addr.String() + "/api/sturdy/" + string(codebaseID) + "/object/" + sha)
		if err != nil {
			return nil, fmt.Errorf("failed to download file from lfs: %w", err)
		}
		defer resp.Body.Close()

		lfsData, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, fmt.Errorf("failed to read file: %w", err)
		}
		return lfsData, nil
	}

	return data, nil
}

var fileExtFilter = map[string]struct{}{
	".jpg":  {},
	".jpeg": {},
	".png":  {},
	".gif":  {},
	".webp": {},
}

func (s *Service) WorkspaceFileType(ctx context.Context, ws *workspaces.Workspace, filePath string, isNew bool) (file.Type, error) {
	fsys, err := live.WorkspaceFS(s.executorProvider, s.snapshotsRepo, ws, isNew)
	if err != nil {
		return file.UnknownType, fmt.Errorf("failed to create fs: %w", err)
	}

	return s.detectFileTypeOnFs(fsys, filePath)
}

func (s *Service) ChangeFileType(ctx context.Context, ch *changes.Change, filePath string, isNew bool) (file.Type, error) {
	fsys, err := live.ChangeFS(s.executorProvider, ch, isNew)
	if err != nil {
		return file.UnknownType, fmt.Errorf("failed to create fs: %w", err)
	}

	return s.detectFileTypeOnFs(fsys, filePath)
}

func (s *Service) detectFileTypeOnFs(fsys fs.FS, filePath string) (file.Type, error) {
	// ext is not in filter, don't validate it
	ext := path.Ext(filePath)
	if _, ok := fileExtFilter[ext]; !ok {
		return file.UnknownType, nil
	}

	fp, err := fsys.Open(filePath)
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

func (s *Service) WorkspaceChecksum(ctx context.Context, ws *workspaces.Workspace, filePath string, isNew bool) (string, error) {
	fsys, err := live.WorkspaceFS(s.executorProvider, s.snapshotsRepo, ws, isNew)
	if err != nil {
		return "", fmt.Errorf("failed to create fs: %w", err)
	}

	fp, err := fsys.Open(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to open file: %w", err)
	}

	return s.checksum(fp)
}

func (s *Service) checksum(fp fs.File) (string, error) {
	sh := sha1.New()
	_, err := io.Copy(sh, fp)
	if err != nil {
		return "", fmt.Errorf("failed to calculate checksum")
	}

	var sum []byte
	sum = sh.Sum(sum)
	return fmt.Sprintf("%x", sum), nil
}
