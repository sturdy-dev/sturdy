package graphql

import "github.com/graph-gophers/graphql-go"

type CodeCommentContextResolver struct {
	*CommentResolver
}

func (r *CodeCommentContextResolver) ID() graphql.ID {
	return graphql.ID(r.comment.ID)
}

func (r *CodeCommentContextResolver) Path() string {
	return r.comment.Path
}

func (r *CodeCommentContextResolver) LineStart() int32 {
	return int32(r.comment.LineStart)
}

func (r *CodeCommentContextResolver) LineEnd() int32 {
	return int32(r.comment.LineEnd)
}

func (r *CodeCommentContextResolver) LineIsNew() bool {
	return r.comment.LineIsNew
}

func (r *CodeCommentContextResolver) Context() string {
	return *r.comment.Context
}

func (r *CodeCommentContextResolver) ContextStartsAtLine() int32 {
	return int32(*r.comment.ContextStartsAtLine)
}
