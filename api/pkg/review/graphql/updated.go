package graphql

import (
	"context"

	"getsturdy.com/api/pkg/auth"
	"getsturdy.com/api/pkg/events"
	eventsv2 "getsturdy.com/api/pkg/events/v2"
	gqlerrors "getsturdy.com/api/pkg/graphql/errors"
	"getsturdy.com/api/pkg/graphql/resolvers"
	"getsturdy.com/api/pkg/review"

	"go.uber.org/zap"
)

func (r *reviewRootResolver) UpdatedReviews(ctx context.Context) (<-chan resolvers.ReviewResolver, error) {
	userID, err := auth.UserID(ctx)
	if err != nil {
		return nil, gqlerrors.Error(err)
	}

	c := make(chan resolvers.ReviewResolver, 100)

	r.eventSubscriber.OnReviewUpdated(ctx, eventsv2.SubscribeUser(userID), func(ctx context.Context, rev *review.Review) error {
		select {
		case <-ctx.Done():
			return events.ErrClientDisconnected
		case c <- &reviewResolver{root: r, rev: rev}:
		default:
			r.logger.Named("updatedReviews").Error("dropped event",
				zap.Stringer("user_id", userID),
				zap.Int("channel_size", len(c)),
			)
		}

		return nil

	})

	return c, nil
}
