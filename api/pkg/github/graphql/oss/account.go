package oss

import (
	"context"

	gqlerrors "mash/pkg/graphql/errors"
	"mash/pkg/graphql/resolvers"
)

type GitHubAccountRootResolver struct{}

func NewGitHubAccountRootResolver() *GitHubAccountRootResolver {
	return &GitHubAccountRootResolver{}
}

func (*GitHubAccountRootResolver) InteralByID(context.Context, string) (resolvers.GitHubAccountResolver, error) {
	return nil, gqlerrors.ErrNotImplemented
}
