package db

import (
	"context"

	"getsturdy.com/api/pkg/change"
)

type Repository interface {
	Get(ctx context.Context, id change.ID) (*change.Change, error)
	ListByIDs(ctx context.Context, ids ...change.ID) ([]*change.Change, error)
	GetByCommitID(ctx context.Context, commitID, codebaseID string) (*change.Change, error)
	Insert(ctx context.Context, ch change.Change) error
	Update(ctx context.Context, ch change.Change) error
}
