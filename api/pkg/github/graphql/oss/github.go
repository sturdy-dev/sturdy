package oss

import (
	"context"

	gqlerrors "getsturdy.com/api/pkg/graphql/errors"
	"getsturdy.com/api/pkg/graphql/resolvers"
)

type gitHubRootResolver struct{}

func NewGitHubRootResolver() resolvers.GitHubRootResolver {
	return &gitHubRootResolver{}
}

func (r *gitHubRootResolver) GitHubRepositories(_ context.Context) ([]resolvers.GitHubRepositoryResolver, error) {
	return nil, gqlerrors.ErrNotImplemented
}

func (r *gitHubRootResolver) SetupGitHubRepository(_ context.Context, _ resolvers.SetupGitHubRepositoryArgs) (resolvers.CodebaseResolver, error) {
	return nil, gqlerrors.ErrNotImplemented
}
