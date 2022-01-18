package resolvers

import (
	"context"

	"github.com/graph-gophers/graphql-go"
)

type ViewRootResolver interface {
	View(ctx context.Context, args ViewArgs) (ViewResolver, error)

	// Internal
	InternalViewsByUser(userID string) ([]ViewResolver, error)
	InternalLastUsedViewByUser(ctx context.Context, codebaseID, userID string) (ViewResolver, error)

	// Mutations
	OpenWorkspaceOnView(ctx context.Context, args OpenViewArgs) (ViewResolver, error)
	CopyWorkspaceToView(ctx context.Context, args CopyViewArgs) (ViewResolver, error)
	RepairView(ctx context.Context, args struct{ ID graphql.ID }) (ViewResolver, error)
	CreateView(ctx context.Context, args CreateViewArgs) (ViewResolver, error)

	// Subscriptions
	UpdatedView(ctx context.Context, args UpdatedViewArgs) (chan ViewResolver, error)
	UpdatedViews(ctx context.Context) (chan ViewResolver, error)
}

type CreateViewArgs struct {
	Input CreateViewInput
}

type CreateViewInput struct {
	WorkspaceID   graphql.ID
	MountPath     string
	MountHostname string
}

type ViewArgs struct {
	ID graphql.ID
}

type OpenViewArgs struct {
	Input OpenWorkspaceOnViewInput
}

type OpenWorkspaceOnViewInput struct {
	ViewID      graphql.ID
	WorkspaceID graphql.ID
}

type CopyViewArgs struct {
	Input CopyWorkspaceOnViewInput
}

type CopyWorkspaceOnViewInput struct {
	ViewID      graphql.ID
	WorkspaceID graphql.ID
}

type UpdatedViewArgs struct {
	ID graphql.ID
}

type ViewResolver interface {
	ID() graphql.ID
	MountPath() string
	ShortMountPath() string
	MountHostname() string
	LastUsedAt() int32
	CreatedAt() int32
	Author(context.Context) (AuthorResolver, error)
	Workspace(ctx context.Context) (WorkspaceResolver, error)
	Status(ctx context.Context) (ViewStatusResolver, error)
	Codebase(ctx context.Context) (CodebaseResolver, error)
	IgnoredPaths(ctx context.Context) ([]string, error)
	SuggestingWorkspace() WorkspaceResolver
}
