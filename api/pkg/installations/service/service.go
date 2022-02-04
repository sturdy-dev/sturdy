package service

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"sync"

	"getsturdy.com/api/pkg/installations"
	"getsturdy.com/api/pkg/installations/db"
	"getsturdy.com/api/pkg/licenses"
	service_organization "getsturdy.com/api/pkg/organization/service"
	"getsturdy.com/api/pkg/version"

	"github.com/google/uuid"
)

var ErrInvalidLicense = errors.New("invalid license")

type Service struct {
	repo                db.Repository
	organizationService *service_organization.Service

	licenseGuard *sync.RWMutex
	// license is the latest license that was retrieved from the licensing server.
	license *licenses.License
}

func New(
	repo db.Repository,
	organizationService *service_organization.Service,
) *Service {
	return &Service{
		repo:                repo,
		organizationService: organizationService,

		licenseGuard: &sync.RWMutex{},
	}
}

func (svc *Service) HasOrganization(ctx context.Context) (bool, error) {
	_, err := svc.organizationService.GetFirst(ctx)
	switch {
	case err == nil:
		return true, nil
	case errors.Is(err, sql.ErrNoRows):
		return false, nil
	default:
		return false, err
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
