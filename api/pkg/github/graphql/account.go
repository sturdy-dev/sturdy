package graphql

import (
	"context"

	gqlerrors "getsturdy.com/api/pkg/graphql/errors"
	"getsturdy.com/api/pkg/graphql/resolvers"
	"getsturdy.com/api/pkg/users"
)

type GitHubAccountRootResolver struct{}

func NewGitHubAccountRootResolver() *GitHubAccountRootResolver {
	return &GitHubAccountRootResolver{}
}

func (*GitHubAccountRootResolver) InteralByID(context.Context, users.ID) (resolvers.GitHubAccountResolver, error) {
	return nil, gqlerrors.ErrNotImplemented
}
