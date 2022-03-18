package db

import (
	"context"
	"database/sql"

	"getsturdy.com/api/pkg/codebases"
	"getsturdy.com/api/pkg/integrations"
)

var _ IntegrationsRepository = &memory{}

type memory struct {
	byCodebaseID map[codebases.ID][]*integrations.Integration
	byID         map[string]*integrations.Integration
}

func NewInMemory() *memory {
	return &memory{
		byCodebaseID: make(map[codebases.ID][]*integrations.Integration),
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

func (m *memory) ListByCodebaseID(ctx context.Context, codebaseID codebases.ID) ([]*integrations.Integration, error) {
	return m.byCodebaseID[codebaseID], nil
}

func (m *memory) Get(ctx context.Context, id string) (*integrations.Integration, error) {
	if v, ok := m.byID[id]; ok {
		return v, nil
	}
	return nil, sql.ErrNoRows
}
