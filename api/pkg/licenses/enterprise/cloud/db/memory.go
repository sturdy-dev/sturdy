package db

import (
	"context"

	"getsturdy.com/api/pkg/licenses"
)

type memory struct {
	byId             map[licenses.ID]*licenses.License
	byOrganizationID map[string][]*licenses.License
	byKey            map[string]*licenses.License
}

func NewMemory() Repository {
	return &memory{
		byId:             make(map[licenses.ID]*licenses.License),
		byOrganizationID: make(map[string][]*licenses.License),
		byKey:            make(map[string]*licenses.License),
	}
}

func (m *memory) Get(_ context.Context, id licenses.ID) (*licenses.License, error) {
	if license, ok := m.byId[id]; ok {
		return license, nil
	}
	return nil, ErrNotFound
}

func (m *memory) Create(_ context.Context, license *licenses.License) error {
	m.byId[license.ID] = license
	m.byOrganizationID[license.OrganizationID] = append(m.byOrganizationID[license.OrganizationID], license)
	m.byKey[license.Key] = license
	return nil
}

func (m *memory) ListByOrganizationID(_ context.Context, organizationID string) ([]*licenses.License, error) {
	if licenses, ok := m.byOrganizationID[organizationID]; ok {
		return licenses, nil
	}
	return nil, ErrNotFound
}

func (m *memory) GetByKey(_ context.Context, key string) (*licenses.License, error) {
	if license, ok := m.byKey[key]; ok {
		return license, nil
	}
	return nil, ErrNotFound
}
