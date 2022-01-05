package resolvers

import (
	"context"

	"github.com/graph-gophers/graphql-go"
)

type FileRootResolver interface {
	InternalFile(ctx context.Context, codebaseID string, pathsWithFallback ...string) (FileOrDirectoryResolver, error)
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
}

type DirectoryResolver interface {
	ID() graphql.ID
	Path() string
	Children(ctx context.Context) []FileOrDirectoryResolver
	Readme(ctx context.Context) (FileResolver, error)
}
