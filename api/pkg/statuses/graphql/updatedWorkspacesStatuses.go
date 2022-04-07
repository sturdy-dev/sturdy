package graphql

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"getsturdy.com/api/pkg/auth"
	"getsturdy.com/api/pkg/events"
	eventsv2 "getsturdy.com/api/pkg/events/v2"
	gqlerrors "getsturdy.com/api/pkg/graphql/errors"
	"getsturdy.com/api/pkg/graphql/resolvers"
	"getsturdy.com/api/pkg/statuses"

	"go.uber.org/zap"
)

func (r *RootResolver) UpdatedWorkspacesStatuses(ctx context.Context, args resolvers.UpdatedWorkspacesStatusesArgs) (<-chan resolvers.WorkspaceStatusResolver, error) {
	userID, err := auth.UserID(ctx)
	if err != nil {
		return nil, gqlerrors.Error(err)
	}

	workspaceIDs := make([]string, 0, len(args.WorkspaceIds))
	for _, id := range args.WorkspaceIds {
		workspaceIDs = append(workspaceIDs, string(id))
	}

	ww, err := r.workspaceService.ListByIDs(ctx, workspaceIDs...)
	if err != nil {
		return nil, gqlerrors.Error(fmt.Errorf("failed to list change commits: %w", err))
	}

	watchWorkspaces := make(map[string]bool, len(ww))
	for _, w := range ww {
		if err := r.authService.CanRead(ctx, w); err != nil {
			return nil, gqlerrors.Error(err)
		}
		watchWorkspaces[w.ID] = true
	}

	c := make(chan resolvers.WorkspaceStatusResolver, 100)
	didErrorOut := false

	r.eventsSubscriber.OnStatusUpdated(ctx, eventsv2.SubscribeUser(userID), func(ctx context.Context, status *statuses.Status) error {
		snapshot, err := r.snapshotsService.GetByCommitSHA(ctx, status.CommitSHA)
		if errors.Is(err, sql.ErrNoRows) {
			return nil
		} else if err != nil {
			return gqlerrors.Error(err)
		} else if snapshot.WorkspaceID == nil {
			return nil
		}

		if !watchWorkspaces[*snapshot.WorkspaceID] {
			return nil
		}

		resolver := &workspaceResolver{
			resolver: &resolver{status: status, root: r},
		}
		select {
		case <-ctx.Done():
			return events.ErrClientDisconnected
		case c <- resolver:
			if didErrorOut {
				didErrorOut = false
			}
			return nil
		default:
			r.logger.Named("updatedChangesStatuses").Error(
				"dropped subscription event",
				zap.Stringer("user_id", userID),
				zap.Stringer("event_type", eventsv2.StatusUpdated),
				zap.Int("channel_size", len(c)),
			)
		}

		return nil
	})

	return c, nil
}
