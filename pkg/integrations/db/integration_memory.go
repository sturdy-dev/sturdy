package db

import (
	"context"
	"database/sql"
	"mash/pkg/integrations"
)

var _ IntegrationsRepository = &memory{}

type memory struct {
	byCodebaseID map[string][]*integrations.Integration
	byID         map[string]*integrations.Integration
}

func NewInMemory() *memory {
	return &memory{
		byCodebaseID: make(map[string][]*integrations.Integration),
		byID:         make(map[string]*integrations.Integration),
	}
}

func (m *memory) Create(ctx context.Context, cfg *integrations.Integration) error {
	m.byCodebaseID[cfg.CodebaseID] = append(m.byCodebaseID[cfg.CodebaseID], cfg)
	m.byID[cfg.ID] = cfg
	return nil
}

func (m *memory) Update(ctx context.Context, cfg *integrations.Integration) error {
	m.byCodebaseID[cfg.CodebaseID] = append(m.byCodebaseID[cfg.CodebaseID], cfg)
	m.byID[cfg.ID] = cfg
	return nil
}

func (m *memory) ListByCodebaseID(ctx context.Context, codebaseID string) ([]*integrations.Integration, error) {
	return m.byCodebaseID[codebaseID], nil
}

func (m *memory) Get(ctx context.Context, id string) (*integrations.Integration, error) {
	if v, ok := m.byID[id]; ok {
		return v, nil
	}
	return nil, sql.ErrNoRows
}
