package graphql

import (
	"context"

	gqlerrors "mash/pkg/graphql/errors"
	"mash/pkg/graphql/resolvers"
	"mash/pkg/suggestions"
	"mash/pkg/unidiff"

	"github.com/graph-gophers/graphql-go"
)

type Resolver struct {
	root *RootResolver

	suggestion *suggestions.Suggestion
}

func (r *Resolver) ID() graphql.ID {
	return graphql.ID(r.suggestion.ID)
}

func (r *Resolver) Author(ctx context.Context) (resolvers.AuthorResolver, error) {
	return r.root.authorResolver.Author(ctx, graphql.ID(r.suggestion.UserID))
}

func (r *Resolver) For(ctx context.Context) (resolvers.WorkspaceResolver, error) {
	ws, err := r.root.workspaceService.GetByID(ctx, r.suggestion.ForWorkspaceID)
	if err != nil {
		return nil, gqlerrors.Error(err)
	}
	return (*r.root.workspaceResolver).InternalWorkspace(ws), nil
}

func (r *Resolver) Workspace(ctx context.Context) (resolvers.WorkspaceResolver, error) {
	ws, err := r.root.workspaceService.GetByID(ctx, r.suggestion.WorkspaceID)
	if err != nil {
		return nil, gqlerrors.Error(err)
	}
	return (*r.root.workspaceResolver).InternalWorkspace(ws), nil
}

func (r *Resolver) DismissedAt() *int32 {
	if r.suggestion.DismissedAt == nil {
		return nil
	}
	dismissedAt := int32(r.suggestion.DismissedAt.Unix())
	return &dismissedAt
}

func (r *Resolver) CreatedAt() int32 {
	return int32(r.suggestion.CreatedAt.Unix())
}

func (r *Resolver) Diffs(ctx context.Context) ([]resolvers.FileDiffResolver, error) {
	allower, err := r.root.authService.GetAllower(ctx, r.suggestion)
	if err != nil {
		return nil, gqlerrors.Error(err)
	}

	diffs, err := r.root.suggestionsService.Diffs(ctx, r.suggestion, unidiff.WithAllower(allower))
	if err != nil {
		return nil, gqlerrors.Error(err)
	}

	rr := make([]resolvers.FileDiffResolver, 0, len(diffs))
	for _, diff := range diffs {
		rr = append(rr, r.root.fileDiffResolver.InternalFileDiff(&diff))
	}
	return rr, nil
}
