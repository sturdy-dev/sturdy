package db

import (
	"context"
	"sort"

	"getsturdy.com/api/pkg/licenses"
	"getsturdy.com/api/pkg/licenses/enterprise/cloud/validations"
)

type memory struct {
	byLicenseID map[licenses.ID][]*validations.Validation
}

func NewMemory() Repository {
	return &memory{
		byLicenseID: make(map[licenses.ID][]*validations.Validation),
	}
}

func (m *memory) Create(_ context.Context, valitions *validations.Validation) error {
	m.byLicenseID[valitions.LicenseID] = append(m.byLicenseID[valitions.LicenseID], valitions)
	return nil
}

func (m *memory) ListLatest(ctx context.Context, licenseID licenses.ID) ([]*validations.Validation, error) {
	list := m.byLicenseID[licenseID]
	sort.Slice(list, func(i, j int) bool {
		return list[i].Timestamp.After(list[j].Timestamp)
	})
	if len(list) <= 10 {
		return list, nil
	}
	return list[:10], nil
}
