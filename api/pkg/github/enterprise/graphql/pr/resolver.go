package pr

import (
	"context"
	"fmt"

	"getsturdy.com/api/pkg/github"
	"getsturdy.com/api/pkg/graphql/resolvers"

	"github.com/graph-gophers/graphql-go"
)

type prResolver struct {
	root *prRootResolver
	pr   *github.PullRequest
}

func (r *prResolver) PullRequestNumber() int32 {
	return int32(r.pr.GitHubPRNumber)
}

func (r *prResolver) Open() bool {
	if r.pr.State == github.PullRequestStateOpen {
		return true
	}
	if r.pr.State == github.PullRequestStateMerging {
		return true
	}
	return false
}

func (r *prResolver) Merged() bool {
	return r.pr.State == github.PullRequestStateMerged
}

func (r *prResolver) State() (resolvers.GitHubPullRequestState, error) {
	switch r.pr.State {
	case github.PullRequestStateOpen:
		return resolvers.GitHubPullRequestStateOpen, nil
	case github.PullRequestStateClosed:
		return resolvers.GitHubPullRequestStateClosed, nil
	case github.PullRequestStateMerged:
		return resolvers.GitHubPullRequestStateMerged, nil
	case github.PullRequestStateMerging:
		return resolvers.GitHubPullRequestStateMerging, nil
	default:
		return "", fmt.Errorf("unknown status: %s", r.pr.State)
	}
}

func (r *prResolver) MergedAt() *int32 {
	if r.pr.MergedAt == nil {
		return nil
	}
	ts := int32(r.pr.MergedAt.Unix())
	return &ts
}

func (r *prResolver) Base() string {
	return r.pr.Base
}

func (r *prResolver) Workspace(ctx context.Context) (resolvers.WorkspaceResolver, error) {
	t := true
	return (*r.root.workspaceResolver).Workspace(ctx, resolvers.WorkspaceArgs{
		ID:            graphql.ID(r.pr.WorkspaceID),
		AllowArchived: &t,
	})
}

func (r *prResolver) ID() graphql.ID {
	return graphql.ID(r.pr.ID)
}

func (r *prResolver) Statuses(ctx context.Context) ([]resolvers.GitHubPullRequestStatusResolver, error) {
	return (*r.root.statusesRootResolver).InternalGitHubPullRequestStatuses(ctx, r.pr)
}
