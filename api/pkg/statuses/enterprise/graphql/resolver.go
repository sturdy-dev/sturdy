package graphql

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"getsturdy.com/api/pkg/auth"
	service_auth "getsturdy.com/api/pkg/auth/service"
	"getsturdy.com/api/pkg/events"
	db_github "getsturdy.com/api/pkg/github/enterprise/db"
	gqlerrors "getsturdy.com/api/pkg/graphql/errors"
	"getsturdy.com/api/pkg/graphql/resolvers"
	"getsturdy.com/api/pkg/statuses/graphql"
	service_statuses "getsturdy.com/api/pkg/statuses/service"

	"go.uber.org/zap"
)

type RootResolver struct {
	*graphql.RootResolver

	logger        *zap.Logger
	statusService *service_statuses.Service
	authService   *service_auth.Service
	githHubPRRepo db_github.GitHubPRRepo
	eventsReader  events.EventReader
}

func New(
	ossResolver *graphql.RootResolver,

	logger *zap.Logger,
	statusService *service_statuses.Service,
	authService *service_auth.Service,
	githHubPRRepo db_github.GitHubPRRepo,
	eventsReader events.EventReader,
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
	didErrorOut := false

	callbackFunc := func(eventType events.EventType, reference string) error {
		if eventType != events.StatusUpdated {
			return nil
		}

		status, err := r.statusService.Get(ctx, reference)
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

		resolver := r.InternalStatus(status)
		select {
		case <-ctx.Done():
			return events.ErrClientDisconnected
		case c <- resolver:
			if didErrorOut {
				didErrorOut = false
			}
			return nil
		default:
			r.logger.Named("updatedGitHubPullRequestStatuses").Error(
				"dropped subscription event",
				zap.Stringer("user_id", userID),
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
