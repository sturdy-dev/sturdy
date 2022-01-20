package db

import (
	"context"
	"database/sql"

	"mash/pkg/codebase"
)

var _ CodebaseRepository = &memory{}

type memory struct {
	byShortID    map[string]*codebase.Codebase
	byInviteCode map[string]*codebase.Codebase
	byID         map[string]*codebase.Codebase
}

func NewMemory() *memory {
	return &memory{
		byShortID:    make(map[string]*codebase.Codebase),
		byInviteCode: make(map[string]*codebase.Codebase),
		byID:         make(map[string]*codebase.Codebase),
	}
}

func (m *memory) Create(entity codebase.Codebase) error {
	m.byID[entity.ID] = &entity
	m.byShortID[string(entity.ShortCodebaseID)] = &entity
	if entity.InviteCode != nil {
		m.byInviteCode[*entity.InviteCode] = &entity
	}
	return nil
}

func (m *memory) Get(id string) (*codebase.Codebase, error) {
	found, ok := m.byID[id]
	if !ok || found.ArchivedAt != nil {
		return nil, sql.ErrNoRows
	}
	return found, nil
}

func (m *memory) GetAllowArchived(id string) (*codebase.Codebase, error) {
	found, ok := m.byInviteCode[id]
	if !ok {
		return nil, sql.ErrNoRows
	}
	return found, nil
}

func (m *memory) GetByInviteCode(inviteCode string) (*codebase.Codebase, error) {
	found, ok := m.byInviteCode[inviteCode]
	if !ok || found.ArchivedAt != nil {
		return nil, sql.ErrNoRows
	}
	return found, nil
}

func (m *memory) GetByShortID(shortID string) (*codebase.Codebase, error) {
	found, ok := m.byShortID[shortID]
	if !ok || found.ArchivedAt != nil {
		return nil, sql.ErrNoRows
	}
	return found, nil
}

func (m *memory) Update(entity *codebase.Codebase) error {
	return m.Create(*entity)
}

func (r *memory) ListByOrganization(_ context.Context, id string) ([]*codebase.Codebase, error) {
	var res []*codebase.Codebase
	for _, cb := range r.byID {
		if cb.OrganizationID != nil && *cb.OrganizationID == id {
			res = append(res, cb)
		}
	}
	return res, nil
}
