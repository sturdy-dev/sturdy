package enterprise

import (
	"context"

	service_auth "getsturdy.com/api/pkg/auth/service"
	service_codebase "getsturdy.com/api/pkg/codebase/service"
	"getsturdy.com/api/pkg/github"
	"getsturdy.com/api/pkg/github/config"
	github_client "getsturdy.com/api/pkg/github/enterprise/client"
	db_github "getsturdy.com/api/pkg/github/enterprise/db"
	service_github "getsturdy.com/api/pkg/github/enterprise/service"
	gqlerrors "getsturdy.com/api/pkg/graphql/errors"
	"getsturdy.com/api/pkg/graphql/resolvers"
	db_snapshots "getsturdy.com/api/pkg/snapshots/db"
	"getsturdy.com/api/pkg/snapshots/snapshotter"
	db_workspace "getsturdy.com/api/pkg/workspace/db"
	"getsturdy.com/api/vcs/executor"

	"github.com/graph-gophers/graphql-go"
	"go.uber.org/zap"
)

type codebaseGitHubIntegrationRootResolver struct {
	gitHubRepositoryRepo   db_github.GitHubRepositoryRepo
	gitHubInstallationRepo db_github.GitHubInstallationRepo
	codebaseService        *service_codebase.Service
	gitExecutorProvider    executor.Provider
	logger                 *zap.Logger
	gitHubAppConfig        config.GitHubAppConfig
	gitHubClientProvider   github_client.ClientProvider
	workspaceReader        db_workspace.WorkspaceReader
	workspaceWriter        db_workspace.WorkspaceWriter
	snapshotter            snapshotter.Snapshotter
	snapshotRepo           db_snapshots.Repository
	authService            *service_auth.Service

	workspaceRootResolver *resolvers.WorkspaceRootResolver
	codebaseRootResolver  *resolvers.CodebaseRootResolver

	gitHubService *service_github.Service
}

func NewCodebaseGitHubIntegrationRootResolver(
	gitHubRepositoryRepo db_github.GitHubRepositoryRepo,
	gitHubInstallationRepo db_github.GitHubInstallationRepo,
	gitExecutorProvider executor.Provider,
	logger *zap.Logger,
	gitHubAppConfig config.GitHubAppConfig,
	gitHubClientProvider github_client.ClientProvider,
	workspaceReader db_workspace.WorkspaceReader,
	workspaceWriter db_workspace.WorkspaceWriter,
	snapshotter snapshotter.Snapshotter,
	snapshotRepo db_snapshots.Repository,
	authService *service_auth.Service,
	codebaseService *service_codebase.Service,

	workspaceRootResolver *resolvers.WorkspaceRootResolver,
	codebaseRootResolver *resolvers.CodebaseRootResolver,

	gitHubService *service_github.Service,
) resolvers.CodebaseGitHubIntegrationRootResolver {
	return &codebaseGitHubIntegrationRootResolver{
		gitHubRepositoryRepo:   gitHubRepositoryRepo,
		gitHubInstallationRepo: gitHubInstallationRepo,
		gitExecutorProvider:    gitExecutorProvider,
		logger:                 logger,
		gitHubAppConfig:        gitHubAppConfig,
		gitHubClientProvider:   gitHubClientProvider,
		workspaceReader:        workspaceReader,
		workspaceWriter:        workspaceWriter,
		snapshotter:            snapshotter,
		snapshotRepo:           snapshotRepo,
		authService:            authService,
		codebaseService:        codebaseService,

		workspaceRootResolver: workspaceRootResolver,
		codebaseRootResolver:  codebaseRootResolver,

		gitHubService: gitHubService,
	}
}

func (r *codebaseGitHubIntegrationRootResolver) InternalGitHubRepositoryByID(id string) (resolvers.CodebaseGitHubIntegrationResolver, error) {
	resolver, err := r.resolveByID(graphql.ID(id))
	if err != nil {
		return nil, gqlerrors.Error(err)
	}
	return resolver, nil
}

func (r *codebaseGitHubIntegrationRootResolver) InternalCodebaseGitHubIntegration(ctx context.Context, codebaseID graphql.ID) (resolvers.CodebaseGitHubIntegrationResolver, error) {
	repo, err := r.resolveByCodebaseID(ctx, codebaseID)
	if err != nil {
		return nil, gqlerrors.Error(err)
	}
	return repo, nil
}

func (r *codebaseGitHubIntegrationRootResolver) UpdateCodebaseGitHubIntegration(ctx context.Context, args resolvers.UpdateCodebaseGitHubIntegrationArgs) (resolvers.CodebaseGitHubIntegrationResolver, error) {
	repo, err := r.gitHubRepositoryRepo.GetByID(string(args.Input.ID))
	if err != nil {
		return nil, gqlerrors.Error(err)
	}

	if err := r.authService.CanWrite(ctx, repo); err != nil {
		return nil, gqlerrors.Error(err)
	}

	if args.Input.Enabled != nil {
		repo.IntegrationEnabled = *args.Input.Enabled
	}
	if args.Input.GitHubIsSourceOfTruth != nil {
		repo.GitHubSourceOfTruth = *args.Input.GitHubIsSourceOfTruth
	}

	err = r.gitHubRepositoryRepo.Update(repo)
	if err != nil {
		return nil, gqlerrors.Error(err)
	}

	resolver, err := r.resolveByID(args.Input.ID)
	if err != nil {
		return nil, gqlerrors.Error(err)
	}

	return resolver, nil
}

func (r *codebaseGitHubIntegrationRootResolver) resolveByCodebaseID(ctx context.Context, codebaseID graphql.ID) (*codebaseGitHubIntegrationResolver, error) {
	repo, err := r.gitHubRepositoryRepo.GetByCodebaseID(string(codebaseID))
	if err != nil {
		return nil, err
	}

	installation, err := r.gitHubInstallationRepo.GetByInstallationID(repo.InstallationID)
	if err != nil {
		return nil, err
	}

	return &codebaseGitHubIntegrationResolver{gitHubRepo: repo, installation: installation, root: r}, nil
}

func (r *codebaseGitHubIntegrationRootResolver) resolveByID(gitHubRepoID graphql.ID) (*codebaseGitHubIntegrationResolver, error) {
	repo, err := r.gitHubRepositoryRepo.GetByID(string(gitHubRepoID))
	if err != nil {
		return nil, err
	}

	installation, err := r.gitHubInstallationRepo.GetByInstallationID(repo.InstallationID)
	if err != nil {
		return nil, err
	}

	return &codebaseGitHubIntegrationResolver{gitHubRepo: repo, installation: installation, root: r}, nil
}

type codebaseGitHubIntegrationResolver struct {
	gitHubRepo   *github.GitHubRepository
	installation *github.GitHubInstallation
	root         *codebaseGitHubIntegrationRootResolver
}

func (r *codebaseGitHubIntegrationResolver) ID() graphql.ID {
	return graphql.ID(r.gitHubRepo.ID)
}

func (r *codebaseGitHubIntegrationResolver) Owner() string {
	return r.installation.Owner
}

func (r *codebaseGitHubIntegrationResolver) Name() string {
	return r.gitHubRepo.Name
}

func (r *codebaseGitHubIntegrationResolver) CreatedAt() int32 {
	return int32(r.gitHubRepo.CreatedAt.Unix())
}

func (r *codebaseGitHubIntegrationResolver) UninstalledAt() *int32 {
	if r.gitHubRepo.UninstalledAt == nil {
		return nil
	}
	t := int32(r.gitHubRepo.UninstalledAt.Unix())
	return &t
}

func (r *codebaseGitHubIntegrationResolver) TrackedBranch() *string {
	if r.gitHubRepo.TrackedBranch == "" {
		return nil
	}
	return &r.gitHubRepo.TrackedBranch
}

func (r *codebaseGitHubIntegrationResolver) SyncedAt() *int32 {
	if r.gitHubRepo.SyncedAt == nil {
		return nil
	}
	t := int32(r.gitHubRepo.SyncedAt.Unix())
	return &t
}

func (r *codebaseGitHubIntegrationResolver) Enabled() bool {
	return r.gitHubRepo.IntegrationEnabled
}

func (r *codebaseGitHubIntegrationResolver) GitHubIsSourceOfTruth() bool {
	return r.gitHubRepo.GitHubSourceOfTruth
}

func (r *codebaseGitHubIntegrationResolver) LastPushErrorMessage() *string {
	return r.gitHubRepo.LastPushErrorMessage
}

func (r *codebaseGitHubIntegrationResolver) LastPushAt() *int32 {
	if r.gitHubRepo.LastPushAt == nil {
		return nil
	}
	t := int32(r.gitHubRepo.LastPushAt.Unix())
	return &t
}

func (r *codebaseGitHubIntegrationResolver) Codebase(ctx context.Context) (resolvers.CodebaseResolver, error) {
	codebaseID := graphql.ID(r.gitHubRepo.CodebaseID)
	return (*r.root.codebaseRootResolver).Codebase(ctx, resolvers.CodebaseArgs{
		ID: &codebaseID,
	})
}
