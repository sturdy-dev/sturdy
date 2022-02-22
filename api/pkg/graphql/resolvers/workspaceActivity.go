package resolvers

import (
	"context"

	"github.com/graph-gophers/graphql-go"
)

type WorkspaceActivityRootResolver interface {
	InternalActivityByWorkspace(ctx context.Context, workspaceID string, args WorkspaceActivityArgs) ([]WorkspaceActivityResolver, error)
	InternalActivityCountByWorkspaceID(context.Context, string) (int32, error)

	ReadWorkspaceActivity(ctx context.Context, args WorkspaceActivityReadArgs) (WorkspaceActivityResolver, error)

	UpdatedWorkspaceActivity(ctx context.Context) (chan WorkspaceActivityResolver, error)
}

type WorkspaceActivityResolver interface {
	ToWorkspaceCommentActivity() (WorkspaceCommentActivityResolver, bool)
	ToWorkspaceCreatedChangeActivity() (WorkspaceCreatedChangeActivityResolver, bool)
	ToWorkspaceRequestedReviewActivity() (WorkspaceRequestedReviewActivityResolver, bool)
	ToWorkspaceReviewedActivity() (WorkspaceReviewedActivityResolver, bool)

	ID() graphql.ID
	Author(context.Context) (AuthorResolver, error)
	CreatedAt() int32
	IsRead(ctx context.Context) (bool, error)
	Workspace(ctx context.Context) (WorkspaceResolver, error)
}

type common interface {
	ID() graphql.ID
	Author(context.Context) (AuthorResolver, error)
	CreatedAt() int32
	IsRead(ctx context.Context) (bool, error)
	Workspace(ctx context.Context) (WorkspaceResolver, error)
}

type WorkspaceCommentActivityResolver interface {
	common
	Comment(context.Context) (CommentResolver, error)
}

type WorkspaceCreatedChangeActivityResolver interface {
	common
	Change(context.Context) (ChangeResolver, error)
}

type WorkspaceRequestedReviewActivityResolver interface {
	common
	Review(context.Context) (ReviewResolver, error)
}

type WorkspaceReviewedActivityResolver interface {
	common
	Review(context.Context) (ReviewResolver, error)
}
