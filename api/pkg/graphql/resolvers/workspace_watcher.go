package resolvers

import (
	"context"

	"getsturdy.com/api/pkg/workspaces"

	"github.com/graph-gophers/graphql-go"
)

type WorkspaceWatcherRootResolver interface {
	// Mutations
	WatchWorkspace(context.Context, WatchWorkspaceArgs) (WorkspaceWatcherResolver, error)
	UnwatchWorkspace(context.Context, UnwatchWorkspaceArgs) (WorkspaceWatcherResolver, error)

	// Subscriptions
	UpdatedWorkspaceWatchers(context.Context, UpdatedWorkspaceWatchersArgs) (<-chan WorkspaceWatcherResolver, error)

	// Internal
	InternalWorkspaceWatchers(context.Context, *workspaces.Workspace) ([]WorkspaceWatcherResolver, error)
}

type UnwatchWorkspaceInput struct {
	WorkspaceID graphql.ID
}

type UnwatchWorkspaceArgs struct {
	Input UnwatchWorkspaceInput
}

type WatchWorkspaceInput struct {
	WorkspaceID graphql.ID
}

type WatchWorkspaceArgs struct {
	Input WatchWorkspaceInput
}

type UpdatedWorkspaceWatchersArgs struct {
	WorkspaceID graphql.ID
}

type WorkspaceWatcherResolver interface {
	User(context.Context) (UserResolver, error)
	Workspace(context.Context) (WorkspaceResolver, error)
	Status() (WorkspaceWatcherStatusType, error)
}

type WorkspaceWatcherStatusType string

const (
	WorkspaceWatcherStatusUndefined WorkspaceWatcherStatusType = ""
	WorkspaceWatcherStatusWatching  WorkspaceWatcherStatusType = "Watching"
	WorkspaceWatcherStatusIgnored   WorkspaceWatcherStatusType = "Ignored"
)
