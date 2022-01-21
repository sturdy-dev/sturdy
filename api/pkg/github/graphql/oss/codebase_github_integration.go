package oss

import (
	"context"

	gqlerrors "getsturdy.com/api/pkg/graphql/errors"
	"getsturdy.com/api/pkg/graphql/resolvers"

	"github.com/graph-gophers/graphql-go"
)

type codebaseGitHubIntegrationRootResolver struct{}

func NewCodebaseGitHubIntegrationRootResolver() resolvers.CodebaseGitHubIntegrationRootResolver {
	return &codebaseGitHubIntegrationRootResolver{}
}

func (r *codebaseGitHubIntegrationRootResolver) InternalGitHubRepositoryByID(id string) (resolvers.CodebaseGitHubIntegrationResolver, error) {
	return nil, gqlerrors.ErrNotImplemented
}

func (r *codebaseGitHubIntegrationRootResolver) InternalCodebaseGitHubIntegration(ctx context.Context, codebaseID graphql.ID) (resolvers.CodebaseGitHubIntegrationResolver, error) {
	return nil, gqlerrors.ErrNotImplemented
}

func (r *codebaseGitHubIntegrationRootResolver) UpdateCodebaseGitHubIntegration(ctx context.Context, args resolvers.UpdateCodebaseGitHubIntegrationArgs) (resolvers.CodebaseGitHubIntegrationResolver, error) {
	return nil, gqlerrors.ErrNotImplemented
}

func (r *codebaseGitHubIntegrationRootResolver) resolveByCodebaseID(ctx context.Context, codebaseID graphql.ID) (resolvers.CodebaseResolver, error) {
	return nil, gqlerrors.ErrNotImplemented
}

func (r *codebaseGitHubIntegrationRootResolver) resolveByID(gitHubRepoID graphql.ID) (resolvers.CodebaseResolver, error) {
	return nil, gqlerrors.ErrNotImplemented
}

func (r *codebaseGitHubIntegrationRootResolver) CreateWorkspaceFromGitHubBranch(ctx context.Context, args resolvers.CreateWorkspaceFromGitHubBranchArgs) (resolvers.WorkspaceResolver, error) {
	return nil, gqlerrors.ErrNotImplemented
}

func (r *codebaseGitHubIntegrationRootResolver) ImportGitHubPullRequests(ctx context.Context, args resolvers.ImportGitHubPullRequestsInputArgs) (resolvers.CodebaseResolver, error) {
	return nil, gqlerrors.ErrNotImplemented
}

func (r *codebaseGitHubIntegrationRootResolver) RefreshGitHubCodebases(ctx context.Context) ([]resolvers.CodebaseResolver, error) {
	return nil, gqlerrors.ErrNotImplemented
}
