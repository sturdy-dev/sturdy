package graphql

import (
	"context"
	"net/url"

	"github.com/graph-gophers/graphql-go"

	"getsturdy.com/api/pkg/file"
	"getsturdy.com/api/pkg/graphql/resolvers"
	"getsturdy.com/api/pkg/workspaces"
)

var _ resolvers.FileInfoResolver = (*fileInfoResolver)(nil)

type fileInfoResolver struct {
	root *fileRootResolver

	id        graphql.ID
	filePath  string
	workspace *workspaces.Workspace
	isNew     bool
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
		if f.isNew {
			q.Set("is_new", "1")
		} else {
			q.Set("is_new", "0")
		}
	} else {
		return nil
	}

	r.RawQuery = q.Encode()
	s := r.String()
	return &s
}

func (f *fileInfoResolver) FileType(ctx context.Context) (resolvers.FileType, error) {
	fileType, err := f.root.fileService.WorkspaceFileType(ctx, f.workspace, f.filePath, f.isNew)
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
