package enterprise

import (
	"context"
	"database/sql"
	"errors"

	"getsturdy.com/api/pkg/analytics"
	"getsturdy.com/api/pkg/auth"
	service_auth "getsturdy.com/api/pkg/auth/service"
	db_codebase "getsturdy.com/api/pkg/codebase/db"
	"getsturdy.com/api/pkg/github/config"
	"getsturdy.com/api/pkg/github/enterprise/client"
	"getsturdy.com/api/pkg/github/enterprise/db"
	service_github "getsturdy.com/api/pkg/github/enterprise/service"
	gqlerrors "getsturdy.com/api/pkg/graphql/errors"
	"getsturdy.com/api/pkg/graphql/resolvers"
	db_user "getsturdy.com/api/pkg/users/db"
	db_view "getsturdy.com/api/pkg/view/db"
	"getsturdy.com/api/pkg/events"
	db_workspace "getsturdy.com/api/pkg/workspace/db"

	"github.com/graph-gophers/graphql-go"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"go.uber.org/zap"
)

var (
	concurrentUpdatedPullRequestConnections = promauto.NewGauge(prometheus.GaugeOpts{
		Name:        "sturdy_graphql_concurrent_subscriptions",
		ConstLabels: prometheus.Labels{"subscription": "updatedGitHubPullRequest"},
	})
)

type prRootResolver struct {
	logger               *zap.Logger
	codebaseResolver     *resolvers.CodebaseRootResolver
	workspaceResolver    *resolvers.WorkspaceRootResolver
	statusesRootResolver *resolvers.StatusesRootResolver

	userRepo     db_user.Repository
	codebaseRepo db_codebase.CodebaseRepository

	gitHubAppConfig config.GitHubAppConfig

	gitHubUserRepo         db.GitHubUserRepo
	workspaceReader        db_workspace.WorkspaceReader
	viewRepo               db_view.Repository
	gitHubPRRepo           db.GitHubPRRepo
	gitHubInstallationRepo db.GitHubInstallationRepo
	gitHubRepositoryRepo   db.GitHubRepositoryRepo

	gitHubClientProvider         client.ClientProvider
	gitHubPersonalClientProvider client.PersonalClientProvider
	events                       events.EventReadWriter
	analyticsClient              analytics.Client

	authService   *service_auth.Service
	gitHubService *service_github.Service
}

func NewResolver(
	logger *zap.Logger,

	codebaseResolver *resolvers.CodebaseRootResolver,
	workspaceResolver *resolvers.WorkspaceRootResolver,
	statusesRootResolver *resolvers.StatusesRootResolver,

	userRepo db_user.Repository,
	codebaseRepo db_codebase.CodebaseRepository,
	workspaceReader db_workspace.WorkspaceReader,
	viewRepo db_view.Repository,

	gitHubAppConfig config.GitHubAppConfig,

	gitHubUserRepo db.GitHubUserRepo,
	gitHubPRRepo db.GitHubPRRepo,
	gitHubInstallationRepo db.GitHubInstallationRepo,
	gitHubRepositoryRepo db.GitHubRepositoryRepo,

	gitHubClientProvider client.ClientProvider,
	gitHubPersonalClientProvider client.PersonalClientProvider,
	events events.EventReadWriter,
	analyticsClient analytics.Client,

	authService *service_auth.Service,
	gitHubService *service_github.Service,
) resolvers.GitHubPullRequestRootResolver {
	return &prRootResolver{
		logger: logger,

		codebaseResolver:     codebaseResolver,
		workspaceResolver:    workspaceResolver,
		statusesRootResolver: statusesRootResolver,

		userRepo:     userRepo,
		codebaseRepo: codebaseRepo,

		gitHubAppConfig: gitHubAppConfig,

		gitHubUserRepo:         gitHubUserRepo,
		workspaceReader:        workspaceReader,
		viewRepo:               viewRepo,
		gitHubPRRepo:           gitHubPRRepo,
		gitHubInstallationRepo: gitHubInstallationRepo,
		gitHubRepositoryRepo:   gitHubRepositoryRepo,

		gitHubClientProvider:         gitHubClientProvider,
		gitHubPersonalClientProvider: gitHubPersonalClientProvider,
		events:                       events,
		analyticsClient:              analyticsClient,

		authService:   authService,
		gitHubService: gitHubService,
	}
}

func (r *prRootResolver) InternalByCodebaseIDAndHeadSHA(ctx context.Context, codebaseID string, commitSHA string) (resolvers.GitHubPullRequestResolver, error) {
	pr, err := r.gitHubPRRepo.GetByCodebaseIDaAndHeadSHA(ctx, codebaseID, commitSHA)
	if err != nil {
		return nil, gqlerrors.Error(err)
	}
	return &prResolver{root: r, pr: pr}, nil
}

// InternalGitHubPullRequest is only to be used in contexts where the request is already authenticated
func (r *prRootResolver) InternalGitHubPullRequestByWorkspaceID(ctx context.Context, args resolvers.GitHubPullRequestArgs) (resolvers.GitHubPullRequestResolver, error) {
	if args.WorkspaceID == nil {
		return nil, gqlerrors.Error(gqlerrors.ErrBadRequest, "workspaceID", "can't be empty")
	}
	prs, err := r.gitHubPRRepo.ListOpenedByWorkspace(string(*args.WorkspaceID))
	if err != nil {
		return nil, gqlerrors.Error(err)
	}
	if len(prs) == 0 {
		// There is no currently open Pull request for this workspace, return the one that was most recently closed
		pr, err := r.gitHubPRRepo.GetMostRecentlyClosedByWorkspace(string(*args.WorkspaceID))
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, nil
		case err == nil:
			return &prResolver{root: r, pr: pr}, nil
		default:
			return nil, gqlerrors.Error(err)
		}

	}
	if len(prs) > 1 {
		r.logger.Warn("more than one opened pull requests for a workspace - this is an erroneous state", zap.Error(err), zap.String("workspace_id", string(*args.WorkspaceID)))
	}
	return &prResolver{root: r, pr: prs[0]}, nil
}

func (r *prRootResolver) CreateOrUpdateGitHubPullRequest(ctx context.Context, args resolvers.CreateOrUpdateGitHubPullRequestArgs) (resolvers.GitHubPullRequestResolver, error) {
	ws, err := r.workspaceReader.Get(string(args.Input.WorkspaceID))
	if err != nil {
		return nil, err
	}

	if err := r.authService.CanWrite(ctx, ws); err != nil {
		return nil, err
	}

	pr, err := r.gitHubService.CreateOrUpdatePullRequest(ctx, ws, args.Input.PatchIDs)
	switch {
	case errors.Is(err, service_github.ErrIntegrationNotEnabled):
		return nil, gqlerrors.Error(gqlerrors.ErrBadRequest, "message", "Pull Requests can only be opened if the integration is enabled and GitHub is considered to be the source of truth")
	case err != nil:
		return nil, gqlerrors.Error(err)
	default:
		return &prResolver{root: r, pr: pr}, nil
	}
}

func (r prRootResolver) UpdatedGitHubPullRequest(ctx context.Context, args resolvers.UpdatedGitHubPullRequestArgs) (<-chan resolvers.GitHubPullRequestResolver, error) {
	userID, err := auth.UserID(ctx)
	if err != nil {
		return nil, gqlerrors.Error(err)
	}

	if args.WorkspaceID == "" {
		return nil, gqlerrors.Error(gqlerrors.ErrBadRequest, "workspaceID", "can't be empty")
	}

	ws, err := r.workspaceReader.Get(string(args.WorkspaceID))
	if err != nil {
		return nil, gqlerrors.Error(err)
	}

	if err := r.authService.CanRead(ctx, ws); err != nil {
		return nil, gqlerrors.Error(err)
	}

	res := make(chan resolvers.GitHubPullRequestResolver, 100)
	didErrorOut := false

	concurrentUpdatedPullRequestConnections.Inc()

	cancelFunc := r.events.SubscribeUser(userID, func(et events.EventType, reference string) error {
		if et != events.GitHubPRUpdated {
			return nil
		}

		if reference != ws.ID {
			return nil
		}

		id := graphql.ID(ws.ID)
		resolver, err := r.InternalGitHubPullRequestByWorkspaceID(ctx, resolvers.GitHubPullRequestArgs{WorkspaceID: &id})
		if err != nil {
			return err
		}
		select {
		case <-ctx.Done():
			return errors.New("disconnected")
		case res <- resolver:
			if didErrorOut {
				didErrorOut = false
			}
			return nil
		default:
			r.logger.Error("dropped subscription event",
				zap.String("user_id", userID),
				zap.String("codebase_id", ws.CodebaseID),
				zap.Stringer("event_type", et),
				zap.Int("channel_size", len(res)),
			)
			didErrorOut = true
			return nil
		}
	})

	go func() {
		<-ctx.Done()
		cancelFunc()
		close(res)
		concurrentUpdatedPullRequestConnections.Dec()
	}()

	return res, nil
}

func (r *prRootResolver) MergeGitHubPullRequest(ctx context.Context, args resolvers.MergeGitHubPullRequestArgs) (resolvers.GitHubPullRequestResolver, error) {
	if args.Input.WorkspaceID == "" {
		return nil, gqlerrors.Error(gqlerrors.ErrBadRequest, "workspaceID", "can't be empty")
	}

	ws, err := r.workspaceReader.Get(string(args.Input.WorkspaceID))
	if err != nil {
		return nil, gqlerrors.Error(err)
	}

	if err := r.authService.CanWrite(ctx, ws); err != nil {
		return nil, gqlerrors.Error(err)
	}

	if err := r.gitHubService.MergePullRequest(ctx, string(args.Input.WorkspaceID)); err != nil {
		var userErr service_github.GitHubUserError
		if errors.As(err, &userErr) {
			return nil, gqlerrors.Error(gqlerrors.ErrBadRequest, "message", userErr.Error())
		}

		return nil, gqlerrors.Error(err)
	}

	return r.InternalGitHubPullRequestByWorkspaceID(ctx, resolvers.GitHubPullRequestArgs{WorkspaceID: &args.Input.WorkspaceID})
}
