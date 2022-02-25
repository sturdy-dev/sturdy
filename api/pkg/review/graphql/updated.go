package graphql

import (
	"context"
	"fmt"

	"getsturdy.com/api/pkg/auth"
	"getsturdy.com/api/pkg/events"
	gqlerrors "getsturdy.com/api/pkg/graphql/errors"
	"getsturdy.com/api/pkg/graphql/resolvers"

	"go.uber.org/zap"
)

func (r *reviewRootResolver) UpdatedReviews(ctx context.Context) (<-chan resolvers.ReviewResolver, error) {
	userID, err := auth.UserID(ctx)
	if err != nil {
		return nil, gqlerrors.Error(err)
	}

	c := make(chan resolvers.ReviewResolver, 100)

	cancelFunc := r.eventsReader.SubscribeUser(userID, func(eventType events.EventType, reviewID string) error {
		if eventType != events.ReviewUpdated {
			return nil
		}

		resolver, err := r.InternalReview(ctx, reviewID)
		if err != nil {
			return fmt.Errorf("failed to get review: %w", err)
		}

		select {
		case <-ctx.Done():
			return nil
		case c <- resolver:
		default:
			r.logger.Named("updatedReviews").Error("dropped event",
				zap.Stringer("user_id", userID),
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
