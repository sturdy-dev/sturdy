package db

import (
	"context"

	"getsturdy.com/api/pkg/changes"
	"getsturdy.com/api/pkg/codebases"
)

type Repository interface {
	Get(ctx context.Context, id changes.ID) (*changes.Change, error)
	ListByIDs(ctx context.Context, ids ...changes.ID) ([]*changes.Change, error)
	GetByCommitID(ctx context.Context, commitID string, codebaseID codebases.ID) (*changes.Change, error)
	Insert(ctx context.Context, ch changes.Change) error
	Update(ctx context.Context, ch changes.Change) error
	GetByParentChangeID(context.Context, changes.ID) (*changes.Change, error)
}
