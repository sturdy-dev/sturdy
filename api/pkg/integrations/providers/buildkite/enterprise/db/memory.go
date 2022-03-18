package db

import (
	"context"
	"database/sql"

	"getsturdy.com/api/pkg/codebases"
	"getsturdy.com/api/pkg/integrations/providers/buildkite"
)

var _ Repository = &memory{}

type memory struct {
	byID            map[string]*buildkite.Config
	byIntegrationID map[string]*buildkite.Config
}

func NewInMemory() *memory {
	m := &memory{
		byID:            make(map[string]*buildkite.Config),
		byIntegrationID: make(map[string]*buildkite.Config),
	}
	return m
}

func (m *memory) Create(ctx context.Context, cfg *buildkite.Config) error {
	m.byID[cfg.ID] = cfg
	m.byIntegrationID[cfg.IntegrationID] = cfg
	return nil
}

func (m *memory) Update(ctx context.Context, cfg *buildkite.Config) error {
	m.byID[cfg.ID] = cfg
	m.byIntegrationID[cfg.IntegrationID] = cfg
	return nil
}

func (m *memory) GetConfigsByCodebaseID(ctx context.Context, codebaseID codebases.ID) ([]*buildkite.Config, error) {
	var res []*buildkite.Config
	for _, v := range m.byID {
		if v.CodebaseID == codebaseID {
			res = append(res, v)
		}
	}
	return res, nil
}

func (m *memory) GetConfigByIntegrationID(ctx context.Context, integrationID string) (*buildkite.Config, error) {
	codebase, found := m.byIntegrationID[integrationID]
	if !found {
		return nil, sql.ErrNoRows
	}
	return codebase, nil
}
