package graphql

import (
	"context"
	"errors"

	gqlerrors "mash/pkg/graphql/errors"
	"mash/pkg/graphql/resolvers"
	"mash/pkg/review"

	"github.com/graph-gophers/graphql-go"
)

type reviewResolver struct {
	root *reviewRootResolver
	rev  *review.Review
}

func (r *reviewResolver) ID() graphql.ID {
	return graphql.ID(r.rev.ID)
}

func (r *reviewResolver) Author(ctx context.Context) (resolvers.AuthorResolver, error) {
	return r.root.authorRootResolver.Author(ctx, graphql.ID(r.rev.UserID))
}

func (r *reviewResolver) Grade() string {
	return string(r.rev.Grade)
}

func (r *reviewResolver) CreatedAt() int32 {
	return int32(r.rev.CreatedAt.Unix())
}

func (r *reviewResolver) DismissedAt() *int32 {
	if r.rev.DismissedAt == nil {
		return nil
	}
	ts := int32(r.rev.DismissedAt.Unix())
	return &ts
}

func (r *reviewResolver) IsReplaced() bool {
	return r.rev.IsReplaced
}

func (r *reviewResolver) RequestedBy(ctx context.Context) (resolvers.AuthorResolver, error) {
	if r.rev.RequestedBy == nil {
		return nil, nil
	}
	return r.root.authorRootResolver.Author(ctx, graphql.ID(*r.rev.RequestedBy))
}

func (r *reviewResolver) Workspace(ctx context.Context) (resolvers.WorkspaceResolver, error) {
	yes := true
	resolver, err := (*r.root.workspaceRootResolver).Workspace(ctx, resolvers.WorkspaceArgs{
		ID:            graphql.ID(r.rev.WorkspaceID),
		AllowArchived: &yes,
	})
	if errors.Is(err, gqlerrors.ErrNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return resolver, err
}
