package graphql

import (
	"context"

	"getsturdy.com/api/pkg/auth"
	"getsturdy.com/api/pkg/events"
	gqlerrors "getsturdy.com/api/pkg/graphql/errors"
	"getsturdy.com/api/pkg/graphql/resolvers"
	"getsturdy.com/api/pkg/view"

	"go.uber.org/zap"
)

func (r *WorkspaceRootResolver) UpdatedWorkspaceDiffs(ctx context.Context, args resolvers.UpdatedWorkspaceDiffsArgs) (<-chan []resolvers.FileDiffResolver, error) {
	// Get codebaseID by the workspaceID
	ws, err := r.workspaceReader.Get(string(args.WorkspaceID))
	if err != nil {
		return nil, gqlerrors.Error(err)
	}

	if err := r.authService.CanRead(ctx, ws); err != nil {
		return nil, gqlerrors.Error(err)
	}

	userID, err := auth.UserID(ctx)
	if err != nil {
		return nil, gqlerrors.Error(err)
	}

	wsResolver := r.InternalWorkspace(ws)

	c := make(chan []resolvers.FileDiffResolver, 100)
	didErrorOut := false

	// TODO: Migrate all to the new subscriber
	cancelFunc := r.viewEvents.SubscribeUser(userID, func(eventType events.EventType, reference string) error {
		select {
		case <-ctx.Done():
			return events.ErrClientDisconnected
		default:
		}

		workspaceUpdated := eventType == events.WorkspaceUpdated && reference == ws.ID
		workspaceSnapshotUpdated := eventType == events.WorkspaceUpdatedSnapshot && reference == ws.ID
		diffsUpdated := workspaceUpdated || workspaceSnapshotUpdated

		if workspaceUpdated {
			// reload the WS, to support streaming diffs if the view ID changes, etc
			ws, err = r.workspaceReader.Get(string(args.WorkspaceID))
			if err != nil {
				return err
			}
		}

		if !diffsUpdated {
			return nil
		}

		diffs, err := wsResolver.Diffs(ctx)
		if err != nil {
			return err
		}

		select {
		case <-ctx.Done():
			return events.ErrClientDisconnected
		case c <- diffs:
			if didErrorOut {
				didErrorOut = false
			}
			return nil
		default:
			r.logger.Error("dropped subscription event",
				zap.Stringer("user_id", userID),
				zap.Stringer("event_type", eventType),
				zap.Int("channel_size", len(c)),
			)
			didErrorOut = true
			return nil
		}
	})

	onViewUpdated := func(ctx context.Context, view *view.View) error {
		// non matching event
		if ws.ViewID == nil || view.ID != *ws.ViewID {
			return nil
		}
		diffs, err := wsResolver.Diffs(ctx)
		if err != nil {
			return err
		}

		select {
		case c <- diffs:
			if didErrorOut {
				didErrorOut = false
			}
		default:
			r.logger.Error("dropped subscription event",
				zap.Stringer("user_id", userID),
				zap.Int("channel_size", len(c)),
			)
			didErrorOut = true
		}
		return nil
	}

	r.eventsSubscriber.User(ctx, userID).OnViewUpdated(ctx, onViewUpdated)

	go func() {
		<-ctx.Done()
		cancelFunc()
	}()

	return c, nil
}
