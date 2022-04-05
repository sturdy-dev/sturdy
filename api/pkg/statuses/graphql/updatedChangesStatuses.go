package graphql

import (
	"context"
	"fmt"

	"getsturdy.com/api/pkg/auth"
	"getsturdy.com/api/pkg/changes"
	"getsturdy.com/api/pkg/codebases"
	"getsturdy.com/api/pkg/events"
	eventsv2 "getsturdy.com/api/pkg/events/v2"
	gqlerrors "getsturdy.com/api/pkg/graphql/errors"
	"getsturdy.com/api/pkg/graphql/resolvers"
	"getsturdy.com/api/pkg/statuses"

	"go.uber.org/zap"
)

func (r *RootResolver) UpdatedChangesStatuses(ctx context.Context, args resolvers.UpdatedChangesStatusesArgs) (<-chan resolvers.StatusResolver, error) {
	userID, err := auth.UserID(ctx)
	if err != nil {
		return nil, gqlerrors.Error(err)
	}

	changeIDs := make([]changes.ID, 0, len(args.ChangeIDs))
	for _, id := range args.ChangeIDs {
		changeIDs = append(changeIDs, changes.ID(id))
	}

	changes, err := r.changeService.ListChanges(ctx, changeIDs...)
	if err != nil {
		return nil, gqlerrors.Error(fmt.Errorf("failed to list change commits: %w", err))
	}

	watchCommit := map[string]bool{}
	codebaseIDs := map[codebases.ID]bool{}
	for _, ch := range changes {
		if err := r.authService.CanRead(ctx, ch); err != nil {
			return nil, gqlerrors.Error(err)
		}
		if ch.CommitID != nil {
			watchCommit[*ch.CommitID] = true
		}
		codebaseIDs[ch.CodebaseID] = true
	}

	c := make(chan resolvers.StatusResolver, 100)
	didErrorOut := false

	r.eventsSubscriber.OnStatusUpdated(ctx, eventsv2.SubscribeUser(userID), func(ctx context.Context, status *statuses.Status) error {
		if !watchCommit[status.CommitSHA] {
			return nil
		}

		resolver := &resolver{status: status, root: r}
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
