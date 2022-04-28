package db

import (
	"context"

	"getsturdy.com/api/pkg/codebases"
	"getsturdy.com/api/pkg/statuses"
)

type Repository interface {
	Create(ctx context.Context, status *statuses.Status) error
	Get(ctx context.Context, id string) (*statuses.Status, error)
	ListByWorkspaceID(ctx context.Context, workspaceID string) ([]*statuses.Status, error)
	ListByCodebaseIDAndCommitID(ctx context.Context, codebaseID codebases.ID, commitID string) ([]*statuses.Status, error)
}
