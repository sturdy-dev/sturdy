package resolvers

import (
	"context"

	"github.com/graph-gophers/graphql-go"

	"getsturdy.com/api/pkg/changes"
	"getsturdy.com/api/pkg/codebases"
	"getsturdy.com/api/pkg/workspaces"
)

type FileRootResolver interface {
	InternalFile(ctx context.Context, codebase *codebases.Codebase, pathsWithFallback ...string) (FileOrDirectoryResolver, error)
	InternalFileInfoInWorkspace(id graphql.ID, filePath string, workspace *workspaces.Workspace, isNew bool) FileInfoResolver
	InternalFileInfoOnChange(id graphql.ID, filePath string, change *changes.Change, isNew bool) FileInfoResolver
}

type FileOrDirectoryResolver interface {
	Path() string
	ToFile() (FileResolver, bool)
	ToDirectory() (DirectoryResolver, bool)
}

type FileResolver interface {
	ID() graphql.ID
	Path() string
	Contents() string
	MimeType() string
	Info() FileInfoResolver
}

type DirectoryResolver interface {
	ID() graphql.ID
	Path() string
	Children(ctx context.Context) ([]FileOrDirectoryResolver, error)
	Readme(ctx context.Context) (FileResolver, error)
}
