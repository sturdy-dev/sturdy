package service

import (
	"context"
	"errors"
	"fmt"
	"sync"

	"getsturdy.com/api/pkg/installations"
	"getsturdy.com/api/pkg/installations/db"
	"getsturdy.com/api/pkg/licenses"
	"getsturdy.com/api/pkg/version"

	"github.com/google/uuid"
)

var ErrInvalidLicense = errors.New("invalid license")

type Service struct {
	repo db.Repository

	licenseGuard *sync.RWMutex
	// license is the latest license that was retrieved from the licensing server.
	license *licenses.License
}

func New(
	repo db.Repository,
) *Service {
	return &Service{
		repo: repo,

		licenseGuard: &sync.RWMutex{},
	}
}

func (svc *Service) UpdateLicense(ctx context.Context, license *licenses.License) error {
	svc.licenseGuard.Lock()
	svc.license = license
	svc.licenseGuard.Unlock()
	return nil
}

// Get returns the global installation object.
func (svc *Service) Get(ctx context.Context) (*installations.Installation, error) {
	ii, err := svc.repo.ListAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to list installations: %w", err)
	}

	svc.licenseGuard.RLock()
	defer svc.licenseGuard.RUnlock()

	switch len(ii) {
	case 0:
		installation := &installations.Installation{
			ID:      uuid.New().String(),
			Type:    version.Type,
			Version: version.Version,
			License: svc.license,
		}
		if err := svc.repo.Create(ctx, installation); err != nil {
			return nil, fmt.Errorf("failed to create installation: %w", err)
		}
		return installation, nil
	case 1:
		installation := ii[0]
		installation.Type = version.Type
		installation.Version = version.Version
		installation.License = svc.license
		return installation, nil
	default:
		return nil, fmt.Errorf("more than one installation found")
	}
}
