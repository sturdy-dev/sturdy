package oss

import (
	"context"

	gqlerrors "getsturdy.com/api/pkg/graphql/errors"
	"getsturdy.com/api/pkg/graphql/resolvers"
)

type GitHubAccountRootResolver struct{}

func NewGitHubAccountRootResolver() *GitHubAccountRootResolver {
	return &GitHubAccountRootResolver{}
}

func (*GitHubAccountRootResolver) InteralByID(context.Context, string) (resolvers.GitHubAccountResolver, error) {
	return nil, gqlerrors.ErrNotImplemented
}
