package db

import (
	"context"
	"database/sql"

	"getsturdy.com/api/pkg/codebases"
)

var _ CodebaseRepository = &memory{}

type memory struct {
	byShortID    map[string]*codebases.Codebase
	byInviteCode map[string]*codebases.Codebase
	byID         map[string]*codebases.Codebase
}

func NewMemory() *memory {
	return &memory{
		byShortID:    make(map[string]*codebases.Codebase),
		byInviteCode: make(map[string]*codebases.Codebase),
		byID:         make(map[string]*codebases.Codebase),
	}
}

func (m *memory) Create(entity codebases.Codebase) error {
	m.byID[entity.ID] = &entity
	m.byShortID[string(entity.ShortCodebaseID)] = &entity
	if entity.InviteCode != nil {
		m.byInviteCode[*entity.InviteCode] = &entity
	}
	return nil
}

func (m *memory) Get(id string) (*codebases.Codebase, error) {
	found, ok := m.byID[id]
	if !ok || found.ArchivedAt != nil {
		return nil, sql.ErrNoRows
	}
	return found, nil
}

func (m *memory) GetAllowArchived(id string) (*codebases.Codebase, error) {
	found, ok := m.byInviteCode[id]
	if !ok {
		return nil, sql.ErrNoRows
	}
	return found, nil
}

func (m *memory) GetByInviteCode(inviteCode string) (*codebases.Codebase, error) {
	found, ok := m.byInviteCode[inviteCode]
	if !ok || found.ArchivedAt != nil {
		return nil, sql.ErrNoRows
	}
	return found, nil
}

func (m *memory) GetByShortID(shortID string) (*codebases.Codebase, error) {
	found, ok := m.byShortID[shortID]
	if !ok || found.ArchivedAt != nil {
		return nil, sql.ErrNoRows
	}
	return found, nil
}

func (m *memory) Update(entity *codebases.Codebase) error {
	return m.Create(*entity)
}

func (r *memory) ListByOrganization(_ context.Context, id string) ([]*codebases.Codebase, error) {
	var res []*codebases.Codebase
	for _, cb := range r.byID {
		if cb.OrganizationID != nil && *cb.OrganizationID == id {
			res = append(res, cb)
		}
	}
	return res, nil
}

func (r *memory) Count(context.Context) (uint64, error) {
	return uint64(len(r.byID)), nil
}
