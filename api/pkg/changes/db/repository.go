package db

import (
	"context"

	"getsturdy.com/api/pkg/changes"
)

type Repository interface {
	Get(ctx context.Context, id changes.ID) (*changes.Change, error)
	ListByIDs(ctx context.Context, ids ...changes.ID) ([]*changes.Change, error)
	GetByCommitID(ctx context.Context, commitID, codebaseID string) (*changes.Change, error)
	Insert(ctx context.Context, ch changes.Change) error
	Update(ctx context.Context, ch changes.Change) error
}
