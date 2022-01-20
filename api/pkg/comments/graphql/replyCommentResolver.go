package graphql

import (
	"context"
	gqlerrors "getsturdy.com/api/pkg/graphql/errors"
	"getsturdy.com/api/pkg/graphql/resolvers"

	"github.com/graph-gophers/graphql-go"
)

type ReplyCommentResolver struct {
	*CommentResolver
}

func (r *ReplyCommentResolver) Parent(ctx context.Context) (resolvers.TopCommentResolver, error) {
	cr, err := r.root.Comment(ctx, resolvers.CommentArgs{
		ID: graphql.ID(*r.comment.ParentComment),
	})
	if err != nil {
		return nil, err
	}
	if tcr, ok := cr.ToTopComment(); ok {
		return tcr, nil
	}
	return nil, gqlerrors.ErrNotFound
}
