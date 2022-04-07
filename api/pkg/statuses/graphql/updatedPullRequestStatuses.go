package graphql

import (
	"context"

	gqlerrors "getsturdy.com/api/pkg/graphql/errors"
	"getsturdy.com/api/pkg/graphql/resolvers"
)

func (r *RootResolver) UpdatedGitHubPullRequestStatuses(ctx context.Context, args resolvers.UpdatedGitHubPullRequestStatusesArgs) (<-chan resolvers.GitHubPullRequestStatusResolver, error) {
	return nil, gqlerrors.ErrNotImplemented
}
