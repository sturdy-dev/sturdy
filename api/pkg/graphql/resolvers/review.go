package resolvers

import (
	"context"

	"github.com/graph-gophers/graphql-go"
)

type ReviewRootResolver interface {
	// Internal
	InternalReview(ctx context.Context, id string) (ReviewResolver, error)
	InternalReviews(ctx context.Context, workspaceID string) ([]ReviewResolver, error)
	InternalDismissAllInWorkspace(ctx context.Context, workspaceID string) error

	// Mutations
	CreateOrUpdateReview(ctx context.Context, args CreateReviewArgs) (ReviewResolver, error)
	DismissReview(ctx context.Context, args DismissReviewArgs) (ReviewResolver, error)
	RequestReview(ctx context.Context, args RequestReviewArgs) (ReviewResolver, error)

	// Subscriptions
	UpdatedReviews(context.Context) (<-chan ReviewResolver, error)
}

type ReviewResolver interface {
	ID() graphql.ID
	Author(context.Context) (AuthorResolver, error)
	Grade() string
	CreatedAt() int32
	DismissedAt() *int32
	IsReplaced() bool
	Workspace(context.Context) (WorkspaceResolver, error)
	RequestedBy(context.Context) (AuthorResolver, error)
}

type CreateReviewArgs struct {
	Input CreateReviewInput
}

type CreateReviewInput struct {
	WorkspaceID graphql.ID
	Grade       string
}

type DismissReviewArgs struct {
	Input DismissReviewInput
}

type DismissReviewInput struct {
	ID graphql.ID
}

type RequestReviewArgs struct {
	Input RequestReviewInput
}

type RequestReviewInput struct {
	WorkspaceID graphql.ID
	UserID      graphql.ID
}
