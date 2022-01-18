package graphql

import (
	"context"
	"database/sql"
	"errors"

	gqlerrors "mash/pkg/graphql/errors"
	"mash/pkg/graphql/resolvers"

	"github.com/graph-gophers/graphql-go"
)

type TopCommentResolver struct {
	*CommentResolver
}

func (r *TopCommentResolver) Workspace(ctx context.Context) (resolvers.WorkspaceResolver, error) {
	if r.comment.WorkspaceID == nil {
		return nil, nil
	}
	res, err := (*r.root.workspaceResolver).Workspace(ctx, resolvers.WorkspaceArgs{ID: graphql.ID(*r.comment.WorkspaceID)})
	// The workspace has been archived since the comment was created
	if errors.Is(err, gqlerrors.ErrNotFound) {
		return nil, nil
	}
	return res, err
}

func (r *TopCommentResolver) Change(ctx context.Context) (resolvers.ChangeResolver, error) {
	if r.comment.ChangeID == nil {
		return nil, nil
	}
	id := graphql.ID(*r.comment.ChangeID)
	return r.root.changeResolver.Change(ctx, resolvers.ChangeArgs{ID: &id})
}

func (r *TopCommentResolver) Replies() ([]resolvers.ReplyCommentResolver, error) {
	replies, err := r.root.commentsRepo.GetByParent(r.comment.ID)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, gqlerrors.Error(err)
	}

	var res []resolvers.ReplyCommentResolver

	for _, reply := range replies {
		resolver, err := r.root.PreFetchedComment(reply)
		if err != nil {
			return nil, gqlerrors.Error(err)
		}
		if rcr, ok := resolver.ToReplyComment(); ok {
			res = append(res, rcr)
			continue
		}
		return nil, gqlerrors.ErrNotFound
	}

	return res, nil
}

func (r *TopCommentResolver) CodeContext() resolvers.CommentCodeContext {
	if r.comment.Path == "" {
		return nil
	}
	return &CodeCommentContextResolver{r.CommentResolver}
}
