package graphql

import (
	"context"
	"fmt"

	"getsturdy.com/api/pkg/auth"
	service_auth "getsturdy.com/api/pkg/auth/service"
	"getsturdy.com/api/pkg/events"
	eventsv2 "getsturdy.com/api/pkg/events/v2"
	db_github "getsturdy.com/api/pkg/github/enterprise/db"
	gqlerrors "getsturdy.com/api/pkg/graphql/errors"
	"getsturdy.com/api/pkg/graphql/resolvers"
	"getsturdy.com/api/pkg/statuses"
	"getsturdy.com/api/pkg/statuses/graphql"
	service_statuses "getsturdy.com/api/pkg/statuses/service"

	"go.uber.org/zap"
)

type RootResolver struct {
	*graphql.RootResolver

	logger        *zap.Logger
	statusService *service_statuses.Service
	authService   *service_auth.Service
	githHubPRRepo db_github.GitHubPRRepository
	eventsReader  *eventsv2.Subscriber
}

func New(
	ossResolver *graphql.RootResolver,

	logger *zap.Logger,
	statusService *service_statuses.Service,
	authService *service_auth.Service,
	githHubPRRepo db_github.GitHubPRRepository,
	eventsReader *eventsv2.Subscriber,
) *RootResolver {
	return &RootResolver{
		RootResolver: ossResolver,

		logger:        logger,
		statusService: statusService,
		authService:   authService,
		githHubPRRepo: githHubPRRepo,
		eventsReader:  eventsReader,
	}
}

func (r *RootResolver) UpdatedGitHubPullRequestStatuses(ctx context.Context, args resolvers.UpdatedGitHubPullRequestStatusesArgs) (<-chan resolvers.StatusResolver, error) {
	userID, err := auth.UserID(ctx)
	if err != nil {
		return nil, gqlerrors.Error(err)
	}

	pr, err := r.githHubPRRepo.Get(string(args.ID))
	if err != nil {
		return nil, gqlerrors.Error(fmt.Errorf("failed to get pr: %w", err))
	}

	if pr.HeadSHA == nil {
		return nil, nil
	}

	if r.authService.CanRead(ctx, pr); err != nil {
		return nil, gqlerrors.Error(err)
	}

	c := make(chan resolvers.StatusResolver, 100)
	r.eventsReader.OnStatusUpdated(ctx, eventsv2.SubscribeUser(userID), func(ctx context.Context, status *statuses.Status) error {
		if status.CommitID != *pr.HeadSHA {
			return nil
		}

		if status.CodebaseID != pr.CodebaseID {
			return nil
		}

		select {
		case <-ctx.Done():
			return events.ErrClientDisconnected
		case c <- r.InternalStatus(status):
			return nil
		default:
			r.logger.Named("updatedGitHubPullRequestStatuses").Error(
				"dropped subscription event",
				zap.Stringer("user_id", userID),
				zap.Stringer("event_type", eventsv2.StatusUpdated),
				zap.Int("channel_size", len(c)),
			)
			return nil
		}
	})
	return c, nil
}
