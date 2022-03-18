package db

import (
	"context"
	"database/sql"

	"getsturdy.com/api/pkg/changes"
	"getsturdy.com/api/pkg/codebases"
)

type inMemoryChangeRepo struct {
	changes map[changes.ID]*changes.Change
}

func NewInMemoryRepo() Repository {
	return &inMemoryChangeRepo{
		changes: make(map[changes.ID]*changes.Change),
	}
}

func (r *inMemoryChangeRepo) Get(_ context.Context, id changes.ID) (*changes.Change, error) {
	if c, ok := r.changes[id]; ok {
		return c, nil
	}
	return nil, sql.ErrNoRows
}

func (r *inMemoryChangeRepo) ListByIDs(_ context.Context, ids ...changes.ID) ([]*changes.Change, error) {
	var res []*changes.Change
	for _, id := range ids {
		if c, ok := r.changes[id]; ok {
			res = append(res, c)
		}
	}
	return res, nil
}

func (r *inMemoryChangeRepo) GetByCommitID(_ context.Context, commitID string, codebaseID codebases.ID) (*changes.Change, error) {
	for _, c := range r.changes {
		if c.CodebaseID == codebaseID && c.CommitID == nil && *c.CommitID == commitID {
			return c, nil
		}
	}
	return nil, sql.ErrNoRows
}

func (r *inMemoryChangeRepo) Insert(_ context.Context, ch changes.Change) error {
	r.changes[ch.ID] = &ch
	return nil
}

func (r *inMemoryChangeRepo) Update(_ context.Context, ch changes.Change) error {
	r.changes[ch.ID] = &ch
	return nil
}

func (r *inMemoryChangeRepo) GetByParentChangeID(_ context.Context, id changes.ID) (*changes.Change, error) {
	for _, change := range r.changes {
		if change.ParentChangeID != nil && *change.ParentChangeID == id {
			return change, nil
		}
	}
	return nil, sql.ErrNoRows
}
