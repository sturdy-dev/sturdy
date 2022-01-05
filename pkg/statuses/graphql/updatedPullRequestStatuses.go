package graphql

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"mash/pkg/auth"
	gqlerrors "mash/pkg/graphql/errors"
	"mash/pkg/graphql/resolvers"
	"mash/pkg/view/events"

	"go.uber.org/zap"
)

func (r *rootResolver) UpdatedGitHubPullRequestStatuses(ctx context.Context, args resolvers.UpdatedGitHubPullRequestStatusesArgs) (<-chan resolvers.StatusResolver, error) {
	userID, err := auth.UserID(ctx)
	if err != nil {
		return nil, gqlerrors.Error(err)
	}

	pr, err := r.gitHubPrRepo.Get(string(args.ID))
	if err != nil {
		return nil, gqlerrors.Error(fmt.Errorf("failed to get pr: %w", err))
	}

	if r.authService.CanRead(ctx, pr); err != nil {
		return nil, gqlerrors.Error(err)
	}

	c := make(chan resolvers.StatusResolver, 100)
	didErrorOut := false

	callbackFunc := func(eventType events.EventType, reference string) error {
		if eventType != events.StatusUpdated {
			return nil
		}

		if pr.HeadSHA == nil {
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

		if status.CommitID != *pr.HeadSHA {
			return nil
		}

		if status.CodebaseID != pr.CodebaseID {
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
			r.logger.Named("updatedGitHubPullRequestStatuses").Error(
				"dropped subscription event",
				zap.String("user_id", userID),
				zap.Stringer("event_type", eventType),
				zap.Int("channel_size", len(c)),
			)
		}

		return nil
	}

	cancelFunc := r.eventsReader.SubscribeUser(userID, callbackFunc)

	go func() {
		<-ctx.Done()
		cancelFunc()
		close(c)
	}()

	return c, nil
}
