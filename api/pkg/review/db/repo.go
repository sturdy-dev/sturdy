package db

import (
	"context"

	"getsturdy.com/api/pkg/review"
	"getsturdy.com/api/pkg/users"
)

type ReviewRepository interface {
	Create(context.Context, review.Review) error
	Update(context.Context, *review.Review) error
	Get(ctx context.Context, id string) (*review.Review, error)
	GetLatestByUserAndWorkspace(ctx context.Context, userID users.ID, workspaceID string) (*review.Review, error)
	ListLatestByWorkspace(ctx context.Context, workspaceID string) ([]*review.Review, error)
	DismissAllInWorkspace(ctx context.Context, workspaceID string) error
}
