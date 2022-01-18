package graphql

import (
	"context"

	gqlerrors "mash/pkg/graphql/errors"
	"mash/pkg/graphql/resolvers"
)

func (r *RootResolver) UpdatedGitHubPullRequestStatuses(ctx context.Context, args resolvers.UpdatedGitHubPullRequestStatusesArgs) (<-chan resolvers.StatusResolver, error) {
	return nil, gqlerrors.ErrNotImplemented
}
