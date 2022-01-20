package graphql

import (
	"context"
	"fmt"

	"getsturdy.com/api/pkg/auth"
	gqlerrors "getsturdy.com/api/pkg/graphql/errors"
	"getsturdy.com/api/pkg/graphql/resolvers"
	"getsturdy.com/api/pkg/view/events"

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

	cancelFunc := r.eventsReader.SubscribeUser(userID, func(eventType events.EventType, reference string) error {
		if eventType != events.WorkspaceWatchingStatusUpdated {
			return nil
		}

		watcher, err := r.workspaceWatcherService.Get(ctx, userID, reference)
		if err != nil {
			return fmt.Errorf("failed to get watcher: %w", err)
		}

		resolver := &watcherResolver{Watcher: watcher, Root: r}
		select {
		case <-ctx.Done():
			return fmt.Errorf("disconnected")
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

	go func() {
		<-ctx.Done()
		cancelFunc()
		close(c)
	}()

	return c, nil
}
