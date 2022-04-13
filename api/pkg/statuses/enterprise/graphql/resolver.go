package graphql

import (
	"context"
	"fmt"

	"getsturdy.com/api/pkg/auth"
	service_auth "getsturdy.com/api/pkg/auth/service"
	"getsturdy.com/api/pkg/events"
	eventsv2 "getsturdy.com/api/pkg/events/v2"
	"getsturdy.com/api/pkg/github"
	db_github "getsturdy.com/api/pkg/github/enterprise/db"
	gqlerrors "getsturdy.com/api/pkg/graphql/errors"
	"getsturdy.com/api/pkg/graphql/resolvers"
	"getsturdy.com/api/pkg/statuses"
	service_statuses "getsturdy.com/api/pkg/statuses/service"

	"go.uber.org/zap"
)

type RootResolver struct {
	resolvers.StatusesRootResolver

	logger           *zap.Logger
	statusService    *service_statuses.Service
	authService      *service_auth.Service
	githHubPRRepo    db_github.GitHubPRRepository
	eventsReader     *eventsv2.Subscriber
	gitHubPrResovler resolvers.GitHubPullRequestRootResolver
}

func New(
	ossResolver resolvers.StatusesRootResolver,

	logger *zap.Logger,
	statusService *service_statuses.Service,
	authService *service_auth.Service,
	githHubPRRepo db_github.GitHubPRRepository,
	eventsReader *eventsv2.Subscriber,
	gitHubPrResovler resolvers.GitHubPullRequestRootResolver,
) resolvers.StatusesRootResolver {
	return &RootResolver{
		StatusesRootResolver: ossResolver,

		logger:           logger,
		statusService:    statusService,
		authService:      authService,
		githHubPRRepo:    githHubPRRepo,
		eventsReader:     eventsReader,
		gitHubPrResovler: gitHubPrResovler,
	}
}

func (r *RootResolver) InternalGitHubPullRequestStatuses(ctx context.Context, pr *github.PullRequest) ([]resolvers.GitHubPullRequestStatusResolver, error) {
	if pr.HeadSHA == nil {
		return nil, nil
	}
	ss, err := r.statusService.List(ctx, pr.CodebaseID, *pr.HeadSHA)
	if err != nil {
		return nil, gqlerrors.Error(err)
	}
	rr := make([]resolvers.GitHubPullRequestStatusResolver, 0, len(ss))
	for _, s := range ss {
		rr = append(rr, &pullRequestResolver{
			StatusResolver: r.InternalStatus(s),
			pr:             pr,
			root:           r,
		})
	}
	return rr, nil

}

type pullRequestResolver struct {
	resolvers.StatusResolver

	root *RootResolver
	pr   *github.PullRequest
}

func (r *pullRequestResolver) GitHubPullRequest(ctx context.Context) (resolvers.GitHubPullRequestResolver, error) {
	return r.root.gitHubPrResovler.InternalPullRequest(r.pr)
}

func (r *RootResolver) UpdatedGitHubPullRequestStatuses(ctx context.Context, args resolvers.UpdatedGitHubPullRequestStatusesArgs) (<-chan resolvers.GitHubPullRequestStatusResolver, error) {
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

	c := make(chan resolvers.GitHubPullRequestStatusResolver, 100)
	r.eventsReader.OnStatusUpdated(ctx, eventsv2.SubscribeUser(userID), func(ctx context.Context, status *statuses.Status) error {
		if status.CommitSHA != *pr.HeadSHA {
			return nil
		}

		if status.CodebaseID != pr.CodebaseID {
			return nil
		}

		pr, err := r.githHubPRRepo.Get(string(args.ID))
		if err != nil {
			return fmt.Errorf("failed to get pr: %w", err)
		}

		resolver := &pullRequestResolver{
			StatusResolver: r.InternalStatus(status),
			pr:             pr,
			root:           r,
		}
		select {
		case <-ctx.Done():
			return events.ErrClientDisconnected
		case c <- resolver:
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
