package resolvers

import (
	"context"

	"getsturdy.com/api/pkg/changes"
	"github.com/graph-gophers/graphql-go"
)

type ActivityRootResolver interface {
	InternalActivityByWorkspace(ctx context.Context, workspaceID string, args ActivityArgs) ([]ActivityResolver, error)
	InternalActivityByChangeID(context.Context, changes.ID, ActivityArgs) ([]ActivityResolver, error)

	ReadWorkspaceActivity(ctx context.Context, args ActivityReadArgs) (ActivityResolver, error)

	UpdatedWorkspaceActivity(ctx context.Context) (chan ActivityResolver, error)
}

type ActivityReadArgs struct {
	Input ActivityReadInput
}

type ActivityReadInput struct {
	ID graphql.ID
}

type ActivityArgs struct {
	Input *ActivityInput
}

type ActivityInput struct {
	UnreadOnly *bool
	Limit      *int32
}

type ActivityResolver interface {
	ToWorkspaceCommentActivity() (CommentActivityResolver, bool)
	ToWorkspaceCreatedChangeActivity() (CreatedChangeActivityResolver, bool)
	ToWorkspaceRequestedReviewActivity() (RequestedReviewActivityResolver, bool)
	ToWorkspaceReviewedActivity() (ReviewedActivityResolver, bool)

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

type CommentActivityResolver interface {
	common
	Comment(context.Context) (CommentResolver, error)
}

type CreatedChangeActivityResolver interface {
	common
	Change(context.Context) (ChangeResolver, error)
}

type RequestedReviewActivityResolver interface {
	common
	Review(context.Context) (ReviewResolver, error)
}

type ReviewedActivityResolver interface {
	common
	Review(context.Context) (ReviewResolver, error)
}
