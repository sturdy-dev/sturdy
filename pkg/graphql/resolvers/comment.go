package resolvers

import (
	"context"
	"mash/pkg/comments"
	"mash/pkg/workspace"

	"github.com/graph-gophers/graphql-go"
)

type CommentRootResolver interface {
	Comment(ctx context.Context, args CommentArgs) (CommentResolver, error)
	InternalWorkspaceComments(workspace *workspace.Workspace) ([]CommentResolver, error)

	// Mutations
	DeleteComment(ctx context.Context, args DeleteCommentArgs) (CommentResolver, error)
	UpdateComment(ctx context.Context, args UpdateCommentArgs) (CommentResolver, error)
	CreateComment(ctx context.Context, args CreateCommentArgs) (CommentResolver, error)

	// Subscriptions
	UpdatedComment(ctx context.Context, args UpdatedCommentArgs) (<-chan CommentResolver, error)

	// Internal
	PreFetchedComment(c comments.Comment) (CommentResolver, error)
}

type CommentArgs struct {
	ID graphql.ID
}

type UpdatedCommentArgs struct {
	WorkspaceID graphql.ID
	ViewID      *graphql.ID
}

type DeleteCommentArgs struct {
	ID graphql.ID
}

type UpdateCommentArgs struct {
	Input UpdateCommentInput
}

type UpdateCommentInput struct {
	ID      graphql.ID
	Message string
}

type CreateCommentArgs struct {
	Input CreateCommentInput
}

type CreateCommentInput struct {
	Message   string
	InReplyTo *graphql.ID

	Path        *string
	LineStart   *int32
	LineEnd     *int32
	LineIsNew   *bool
	ChangeID    *graphql.ID
	WorkspaceID *graphql.ID
	ViewID      *graphql.ID
}

type CommentResolver interface {
	ToTopComment() (TopCommentResolver, bool)
	ToReplyComment() (ReplyCommentResolver, bool)

	ID() graphql.ID
	Author(context.Context) (AuthorResolver, error)
	CreatedAt() int32
	DeletedAt() *int32
	Message() string
}

type TopCommentResolver interface {
	ID() graphql.ID
	Author(context.Context) (AuthorResolver, error)
	CreatedAt() int32
	DeletedAt() *int32
	Message() string
	Workspace(ctx context.Context) (WorkspaceResolver, error)
	Change(ctx context.Context) (ChangeResolver, error)
	Replies() ([]ReplyCommentResolver, error)
	CodeContext() CommentCodeContext
}

type ReplyCommentResolver interface {
	ID() graphql.ID
	Author(context.Context) (AuthorResolver, error)
	CreatedAt() int32
	DeletedAt() *int32
	Message() string
	Parent(context.Context) (TopCommentResolver, error)
}

type CommentCodeContext interface {
	ID() graphql.ID
	Path() string
	LineStart() int32
	LineEnd() int32
	LineIsNew() bool
	Context() string
	ContextStartsAtLine() int32
}
