package resolvers

import (
	"context"
	"github.com/graph-gophers/graphql-go"
)

type RebaseStatusRootResolver interface {
	InternalWorkspaceRebaseStatus(ctx context.Context, workspaceID string) (RebaseStatusResolver, error)
}

type RebaseStatusResolver interface {
	ID() graphql.ID
	IsRebasing() bool
	ConflictingFiles() ([]ConflictingFileResolver, error)
}
