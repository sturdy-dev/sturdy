package graphql

import (
	"context"

	"getsturdy.com/api/pkg/auth"
	"getsturdy.com/api/pkg/events"
	gqlerrors "getsturdy.com/api/pkg/graphql/errors"
	"getsturdy.com/api/pkg/graphql/resolvers"

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

	// TODO: Support reloading the workspace.ViewID

	cancelFunc := r.viewEvents.SubscribeUser(userID, func(eventType events.EventType, reference string) error {
		workspaceUpdated := eventType == events.WorkspaceUpdated && reference == ws.ID
		viewUpdated := eventType == events.ViewUpdated && ws.ViewID != nil && reference == *ws.ViewID
		workspaceSnapshotUpdated := eventType == events.WorkspaceUpdatedSnapshot && reference == ws.ID
		diffsUpdated := workspaceUpdated || viewUpdated || workspaceSnapshotUpdated

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

	go func() {
		<-ctx.Done()
		cancelFunc()
		close(c)
	}()

	return c, nil
}
