package inmemory

import (
	"database/sql"
	"mash/pkg/codebase"
	db_codebase "mash/pkg/codebase/db"
)

type inMemoryCodebaseRepository struct {
	codebases []codebase.Codebase
}

func NewInMemoryCodebaseRepo() db_codebase.CodebaseRepository {
	return &inMemoryCodebaseRepository{codebases: make([]codebase.Codebase, 0)}
}

func (r *inMemoryCodebaseRepository) Create(entity codebase.Codebase) error {
	r.codebases = append(r.codebases, entity)
	return nil
}

func (r *inMemoryCodebaseRepository) Get(id string) (*codebase.Codebase, error) {
	for _, cb := range r.codebases {
		if cb.ID == id && cb.ArchivedAt == nil {
			return &cb, nil
		}
	}
	return nil, sql.ErrNoRows
}

func (r *inMemoryCodebaseRepository) GetAllowArchived(id string) (*codebase.Codebase, error) {
	for _, cb := range r.codebases {
		if cb.ID == id {
			return &cb, nil
		}
	}
	return nil, sql.ErrNoRows
}

func (r *inMemoryCodebaseRepository) GetByInviteCode(inviteCode string) (*codebase.Codebase, error) {
	for _, cb := range r.codebases {
		if cb.InviteCode == nil && *cb.InviteCode == inviteCode && cb.ArchivedAt == nil {
			return &cb, nil
		}
	}
	return nil, sql.ErrNoRows
}

func (r *inMemoryCodebaseRepository) GetByShortID(shortID string) (*codebase.Codebase, error) {
	for _, cb := range r.codebases {
		if string(cb.ShortCodebaseID) == shortID && cb.ArchivedAt == nil {
			return &cb, nil
		}
	}
	return nil, sql.ErrNoRows
}

func (r *inMemoryCodebaseRepository) Update(entity *codebase.Codebase) error {
	for k, cb := range r.codebases {
		if cb.ID == entity.ID {
			r.codebases[k] = *entity
			return nil
		}
	}
	return sql.ErrNoRows
}
