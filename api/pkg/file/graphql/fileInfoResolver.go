package graphql

import (
	"context"
	"net/url"

	"github.com/graph-gophers/graphql-go"

	"getsturdy.com/api/pkg/changes"
	"getsturdy.com/api/pkg/file"
	"getsturdy.com/api/pkg/graphql/resolvers"
	"getsturdy.com/api/pkg/workspaces"
)

var _ resolvers.FileInfoResolver = (*fileInfoResolver)(nil)

type fileInfoResolver struct {
	root *fileRootResolver

	id       graphql.ID
	filePath string
	isNew    bool

	workspace *workspaces.Workspace
	change    *changes.Change
}

func (f *fileInfoResolver) ID() graphql.ID {
	return f.id
}

func (f *fileInfoResolver) RawURL() *string {
	if f.filePath == "/dev/null" {
		return nil
	}

	var r url.URL
	r.Path = "/v3/file"
	q := r.Query()
	q.Set("path", f.filePath)

	if f.workspace != nil {
		q.Set("workspace_id", f.workspace.ID)
	} else if f.change != nil {
		q.Set("change_id", string(f.change.ID))
	} else {
		return nil
	}

	if f.isNew {
		q.Set("is_new", "1")
	} else {
		q.Set("is_new", "0")
	}

	r.RawQuery = q.Encode()
	s := r.String()
	return &s
}

func (f *fileInfoResolver) FileType(ctx context.Context) (resolvers.FileType, error) {
	var err error
	var fileType file.Type

	if f.workspace != nil {
		fileType, err = f.root.fileService.WorkspaceFileType(ctx, f.workspace, f.filePath, f.isNew)
	} else if f.change != nil {
		fileType, err = f.root.fileService.ChangeFileType(ctx, f.change, f.filePath, f.isNew)
	} else {
		return resolvers.FileTypeBinary, nil
	}
	if err != nil {
		return resolvers.FileTypeBinary, nil
	}

	switch fileType {
	case file.ImageType:
		return resolvers.FileTypeImage, nil
	case file.TextType:
		return resolvers.FileTypeText, nil
	case file.BinaryType:
		return resolvers.FileTypeBinary, nil
	default:
		return resolvers.FileTypeBinary, nil
	}
}
