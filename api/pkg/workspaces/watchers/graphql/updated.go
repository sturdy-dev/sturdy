package graphql

import (
	"context"

	"getsturdy.com/api/pkg/auth"
	events "getsturdy.com/api/pkg/events"
	eventsv2 "getsturdy.com/api/pkg/events/v2"
	gqlerrors "getsturdy.com/api/pkg/graphql/errors"
	"getsturdy.com/api/pkg/graphql/resolvers"
	"getsturdy.com/api/pkg/workspaces/watchers"

	"go.uber.org/zap"
)

func (r *rootResolver) UpdatedWorkspaceWatchers(ctx context.Context, args resolvers.UpdatedWorkspaceWatchersArgs) (<-chan resolvers.WorkspaceWatcherResolver, error) {
	userID, err := auth.UserID(ctx)
	if err != nil {
		return nil, gqlerrors.Error(err)
	}

	workspace, err := r.workspaceService.GetByID(ctx, string(args.WorkspaceID))
	if err != nil {
		return nil, gqlerrors.Error(err)
	}

	if r.authService.CanRead(ctx, workspace); err != nil {
		return nil, gqlerrors.Error(err)
	}

	c := make(chan resolvers.WorkspaceWatcherResolver, 100)
	didErorrOut := false

	r.eventsReader.OnWorkspaceWatchingStatusUpdated(ctx, eventsv2.SubscribeUser(userID), func(ctx context.Context, watcher *watchers.Watcher) error {
		resolver := &watcherResolver{Watcher: watcher, Root: r}
		select {
		case <-ctx.Done():
			return events.ErrClientDisconnected
		case c <- resolver:
			if didErorrOut {
				didErorrOut = false
			}
			return nil
		default:
			r.logger.Error("dropped subscription event", zap.Error(err))
		}

		return nil

	})
	return c, nil
}
