package db

import (
	"context"

	"getsturdy.com/api/pkg/installations/statistics"
)

type memory struct {
	byLicenseKey map[string]*statistics.Statistic
}

func NewMemory() Repository {
	return &memory{}
}

func (m *memory) Create(_ context.Context, statistic *statistics.Statistic) error {
	if statistic.LicenseKey != nil {
		m.byLicenseKey[*statistic.LicenseKey] = statistic
	}
	return nil
}

func (m *memory) GetByLicenseKey(_ context.Context, key string) (*statistics.Statistic, error) {
	statistic, ok := m.byLicenseKey[key]
	if !ok {
		return nil, ErrNotFound
	}
	return statistic, nil
}
