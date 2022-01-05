package graphql

import (
	"context"
	"fmt"

	gqlerrors "mash/pkg/graphql/errors"
	"mash/pkg/graphql/resolvers"
	"mash/pkg/workspace/watchers"

	"github.com/graph-gophers/graphql-go"
)

type watcherResolver struct {
	Root    *rootResolver
	Watcher *watchers.Watcher
}

func (r *watcherResolver) User(ctx context.Context) (resolvers.UserResolver, error) {
	return (*r.Root.userRootResover).InternalUser(ctx, r.Watcher.UserID)
}

func (r *watcherResolver) Workspace(ctx context.Context) (resolvers.WorkspaceResolver, error) {
	return (*r.Root.workspaceRootResolver).Workspace(ctx, resolvers.WorkspaceArgs{
		ID: graphql.ID(r.Watcher.WorkspaceID),
	})
}

func (r *watcherResolver) Status() (resolvers.WorkspaceWatcherStatusType, error) {
	switch r.Watcher.Status {
	case watchers.StatusIgnored:
		return resolvers.WorkspaceWatcherStatusIgnored, nil
	case watchers.StatusWatching:
		return resolvers.WorkspaceWatcherStatusWatching, nil
	default:
		return resolvers.WorkspaceWatcherStatusUndefined, gqlerrors.Error(fmt.Errorf("unknown status %s", r.Watcher.Status))
	}
}
