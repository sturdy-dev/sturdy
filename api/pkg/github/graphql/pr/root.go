package pr

import (
	"context"

	"getsturdy.com/api/pkg/codebases"
	gqlerrors "getsturdy.com/api/pkg/graphql/errors"
	"getsturdy.com/api/pkg/graphql/resolvers"
)

type prRootResolver struct{}

func NewResolver() resolvers.GitHubPullRequestRootResolver {
	return &prRootResolver{}
}

func (r *prRootResolver) InternalByCodebaseIDAndHeadSHA(context.Context, codebases.ID, string) (resolvers.GitHubPullRequestResolver, error) {
	return nil, gqlerrors.ErrNotImplemented
}

func (r *prRootResolver) InternalGitHubPullRequestByWorkspaceID(ctx context.Context, args resolvers.GitHubPullRequestArgs) (resolvers.GitHubPullRequestResolver, error) {
	return nil, gqlerrors.ErrNotImplemented
}

func (r *prRootResolver) CreateOrUpdateGitHubPullRequest(ctx context.Context, args resolvers.CreateOrUpdateGitHubPullRequestArgs) (resolvers.GitHubPullRequestResolver, error) {
	return nil, gqlerrors.ErrNotImplemented
}

func (r prRootResolver) UpdatedGitHubPullRequest(ctx context.Context, args resolvers.UpdatedGitHubPullRequestArgs) (<-chan resolvers.GitHubPullRequestResolver, error) {
	return nil, gqlerrors.ErrNotImplemented
}

func (r *prRootResolver) MergeGitHubPullRequest(ctx context.Context, args resolvers.MergeGitHubPullRequestArgs) (resolvers.GitHubPullRequestResolver, error) {
	return nil, gqlerrors.ErrNotImplemented
}
