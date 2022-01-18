package enterprise

import (
	"context"

	"mash/pkg/github"
	gqlerrors "mash/pkg/graphql/errors"
	"mash/pkg/graphql/resolvers"

	"github.com/graph-gophers/graphql-go"
)

type prResolver struct {
	root *prRootResolver
	pr   *github.GitHubPullRequest
}

func (r *prResolver) PullRequestNumber() int32 {
	return int32(r.pr.GitHubPRNumber)
}

func (r *prResolver) Open() bool {
	return r.pr.Open
}

func (r *prResolver) Merged() bool {
	return r.pr.Merged
}

func (r *prResolver) MergedAt() *int32 {
	if !r.pr.Merged {
		return nil
	}
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
	return (*r.root.workspaceResolver).Workspace(ctx, resolvers.WorkspaceArgs{ID: graphql.ID(r.pr.WorkspaceID)})
}

func (r *prResolver) ID() graphql.ID {
	return graphql.ID(r.pr.ID)
}

func (r *prResolver) Statuses(ctx context.Context) ([]resolvers.StatusResolver, error) {
	if r.pr.HeadSHA == nil {
		return nil, nil
	}

	ws, err := r.root.workspaceReader.Get(r.pr.WorkspaceID)
	if err != nil {
		return nil, gqlerrors.Error(err)
	}

	return (*r.root.statusesRootResolver).InteralStatusesByCodebaseIDAndCommitID(ctx, ws.CodebaseID, *r.pr.HeadSHA)
}
