package inmemory

import (
	"context"
	"database/sql"

	"getsturdy.com/api/pkg/change"
	db_change "getsturdy.com/api/pkg/change/db"
)

type inMemoryChangeRepo struct {
	changes map[change.ID]change.Change
}

func NewInMemoryChangeRepo() db_change.Repository {
	return &inMemoryChangeRepo{
		changes: make(map[change.ID]change.Change),
	}
}

func (r *inMemoryChangeRepo) Get(_ context.Context, id change.ID) (*change.Change, error) {
	if c, ok := r.changes[id]; ok {
		return &c, nil
	}
	return nil, sql.ErrNoRows
}

func (r *inMemoryChangeRepo) ListByIDs(_ context.Context, ids ...change.ID) ([]*change.Change, error) {
	var res []*change.Change
	for _, id := range ids {
		if c, ok := r.changes[id]; ok {
			res = append(res, &c)
		}
	}
	return res, nil
}

func (r *inMemoryChangeRepo) GetByCommitID(_ context.Context, commitID, codebaseID string) (*change.Change, error) {
	for _, c := range r.changes {
		if c.CodebaseID == codebaseID && c.CommitID == nil && *c.CommitID == commitID {
			return &c, nil
		}
	}
	return nil, sql.ErrNoRows
}

func (r *inMemoryChangeRepo) Insert(_ context.Context, ch change.Change) error {
	r.changes[ch.ID] = ch
	return nil
}

func (r *inMemoryChangeRepo) Update(_ context.Context, ch change.Change) error {
	r.changes[ch.ID] = ch
	return nil
}
