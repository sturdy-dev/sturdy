package inmemory

import (
	"context"
	"database/sql"

	"getsturdy.com/api/pkg/codebases"
	db_codebases "getsturdy.com/api/pkg/codebases/db"
)

type inMemoryCodebaseRepository struct {
	codebases []codebases.Codebase
}

func NewInMemoryCodebaseRepo() db_codebases.CodebaseRepository {
	return &inMemoryCodebaseRepository{codebases: make([]codebases.Codebase, 0)}
}

func (r *inMemoryCodebaseRepository) Create(entity codebases.Codebase) error {
	r.codebases = append(r.codebases, entity)
	return nil
}

func (r *inMemoryCodebaseRepository) Get(id codebases.ID) (*codebases.Codebase, error) {
	for _, cb := range r.codebases {
		if cb.ID == id && cb.ArchivedAt == nil {
			return &cb, nil
		}
	}
	return nil, sql.ErrNoRows
}

func (r *inMemoryCodebaseRepository) GetAllowArchived(id codebases.ID) (*codebases.Codebase, error) {
	for _, cb := range r.codebases {
		if cb.ID == id {
			return &cb, nil
		}
	}
	return nil, sql.ErrNoRows
}

func (r *inMemoryCodebaseRepository) GetByInviteCode(inviteCode string) (*codebases.Codebase, error) {
	for _, cb := range r.codebases {
		if cb.InviteCode == nil && *cb.InviteCode == inviteCode && cb.ArchivedAt == nil {
			return &cb, nil
		}
	}
	return nil, sql.ErrNoRows
}

func (r *inMemoryCodebaseRepository) GetByShortID(shortID codebases.ShortCodebaseID) (*codebases.Codebase, error) {
	for _, cb := range r.codebases {
		if cb.ShortCodebaseID == shortID && cb.ArchivedAt == nil {
			return &cb, nil
		}
	}
	return nil, sql.ErrNoRows
}

func (r *inMemoryCodebaseRepository) Update(entity *codebases.Codebase) error {
	for k, cb := range r.codebases {
		if cb.ID == entity.ID {
			r.codebases[k] = *entity
			return nil
		}
	}
	return sql.ErrNoRows
}

func (r *inMemoryCodebaseRepository) ListByOrganization(_ context.Context, id string) ([]*codebases.Codebase, error) {
	var res []*codebases.Codebase
	for _, cb := range r.codebases {
		if cb.OrganizationID != nil && *cb.OrganizationID == id {
			c2 := cb
			res = append(res, &c2)
		}
	}
	return res, nil
}

func (r *inMemoryCodebaseRepository) Count(_ context.Context) (uint64, error) {
	return uint64(len(r.codebases)), nil
}
