package graphql

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"getsturdy.com/api/pkg/auth"
	"getsturdy.com/api/pkg/change"
	gqlerrors "getsturdy.com/api/pkg/graphql/errors"
	"getsturdy.com/api/pkg/graphql/resolvers"
	"getsturdy.com/api/pkg/view/events"

	"go.uber.org/zap"
)

func (r *RootResolver) UpdatedChangesStatuses(ctx context.Context, args resolvers.UpdatedChangesStatusesArgs) (<-chan resolvers.StatusResolver, error) {
	userID, err := auth.UserID(ctx)
	if err != nil {
		return nil, gqlerrors.Error(err)
	}

	changeIDs := make([]change.ID, 0, len(args.ChangeIDs))
	for _, id := range args.ChangeIDs {
		changeIDs = append(changeIDs, change.ID(id))
	}

	changes, err := r.changeService.ListChangeCommits(ctx, changeIDs...)
	if err != nil {
		return nil, gqlerrors.Error(fmt.Errorf("failed to list change commits: %w", err))
	}

	watchCommit := map[string]bool{}
	codebaseIDs := map[string]bool{}
	for _, change := range changes {
		if err := r.authService.CanRead(ctx, change); err != nil {
			return nil, gqlerrors.Error(err)
		}
		watchCommit[change.CommitID] = true
		codebaseIDs[change.CodebaseID] = true
	}

	c := make(chan resolvers.StatusResolver, 100)
	didErrorOut := false

	cancelFunc := r.eventsReader.SubscribeUser(userID, func(eventType events.EventType, reference string) error {
		if eventType != events.StatusUpdated {
			return nil
		}

		status, err := r.svc.Get(ctx, reference)
		switch {
		case err == nil:
		case errors.Is(err, sql.ErrNoRows):
			return nil
		default:
			return fmt.Errorf("failed to get status by id: %w", err)
		}

		if !watchCommit[status.CommitID] {
			return nil
		}

		resolver := &resolver{status: status, root: r}
		select {
		case <-ctx.Done():
			return errors.New("disconneted")
		case c <- resolver:
			if didErrorOut {
				didErrorOut = false
			}
			return nil
		default:
			r.logger.Named("updatedChangesStatuses").Error(
				"dropped subscription event",
				zap.String("user_id", userID),
				zap.Stringer("event_type", eventType),
				zap.Int("channel_size", len(c)),
			)
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
