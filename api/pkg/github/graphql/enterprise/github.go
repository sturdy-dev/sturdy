package enterprise

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strconv"

	"github.com/graph-gophers/graphql-go"

	"getsturdy.com/api/pkg/auth"
	service_github "getsturdy.com/api/pkg/github/enterprise/service"
	gqlerrors "getsturdy.com/api/pkg/graphql/errors"
	"getsturdy.com/api/pkg/graphql/resolvers"
)

type gitHubRootResolver struct {
	svc                  *service_github.Service
	codebaseRootResolver resolvers.CodebaseRootResolver
}

func NewGitHubRootResolver(
	svc *service_github.Service,
	codebaseRootResolver resolvers.CodebaseRootResolver,
) resolvers.GitHubRootResolver {
	return &gitHubRootResolver{
		svc:                  svc,
		codebaseRootResolver: codebaseRootResolver,
	}
}

func (r *gitHubRootResolver) GitHubRepositories(ctx context.Context) ([]resolvers.GitHubRepositoryResolver, error) {
	userID, err := auth.UserID(ctx)
	if err != nil {
		return nil, gqlerrors.Error(err)
	}

	repos, err := r.svc.ListAllAccessibleRepositoriesFromGitHub(userID)
	switch {
	case errors.Is(err, sql.ErrNoRows):
		return nil, nil
	case err != nil:
		return nil, gqlerrors.Error(err)
	}

	var res []resolvers.GitHubRepositoryResolver
	for _, repo := range repos {
		res = append(res, &gitHubRepositoryResolver{root: r, repo: repo})
	}

	return res, nil
}

func (r *gitHubRootResolver) SetupGitHubRepository(ctx context.Context, args resolvers.SetupGitHubRepositoryArgs) (resolvers.CodebaseResolver, error) {
	installationID, err := strconv.ParseInt(string(args.Input.GitHubInstallationID), 10, 64)
	if err != nil {
		return nil, gqlerrors.Error(err)
	}

	repositoryID, _ := strconv.ParseInt(string(args.Input.GitHubRepositoryID), 10, 64)
	if err != nil {
		return nil, gqlerrors.Error(err)
	}

	userID, err := auth.UserID(ctx)
	if err != nil {
		return nil, gqlerrors.Error(err)
	}

	organizationID := string(args.Input.OrganizationID)

	codebase, err := r.svc.CreateNonReadyCodebaseAndCloneByIDs(ctx, installationID, repositoryID, userID, organizationID)
	if err != nil {
		return nil, gqlerrors.Error(err)
	}

	id := graphql.ID(codebase.ID)
	return r.codebaseRootResolver.Codebase(ctx, resolvers.CodebaseArgs{ID: &id})
}

type gitHubRepositoryResolver struct {
	root *gitHubRootResolver
	repo service_github.GitHubRepo
}

func (r *gitHubRepositoryResolver) ID() graphql.ID {
	return graphql.ID(fmt.Sprintf("%d", r.repo.RepositoryID))
}

func (r *gitHubRepositoryResolver) GitHubInstallationID() graphql.ID {
	return graphql.ID(fmt.Sprintf("%d", r.repo.InstallationID))
}

func (r *gitHubRepositoryResolver) GitHubRepositoryID() graphql.ID {
	return graphql.ID(fmt.Sprintf("%d", r.repo.RepositoryID))
}

func (r *gitHubRepositoryResolver) GitHubOwner() string {
	return r.repo.Owner
}

func (r *gitHubRepositoryResolver) GitHubName() string {
	return r.repo.Name
}

func (r *gitHubRepositoryResolver) Codebase(ctx context.Context) (resolvers.CodebaseResolver, error) {
	ghRepo, err := r.root.svc.GetRepositoryByInstallationAndRepoID(ctx, r.repo.InstallationID, r.repo.RepositoryID)
	switch {
	case errors.Is(err, sql.ErrNoRows):
		return nil, nil
	case err != nil:
		return nil, gqlerrors.Error(err)
	}
	id := graphql.ID(ghRepo.CodebaseID)
	return r.root.codebaseRootResolver.Codebase(ctx, resolvers.CodebaseArgs{ID: &id})
}
