package service

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"sync"
	"time"

	"getsturdy.com/api/pkg/installations"
	"getsturdy.com/api/pkg/installations/db"
	"getsturdy.com/api/pkg/licenses"
	validator_license "getsturdy.com/api/pkg/licenses/enterprise/selfhosted/validator"
	service_organization "getsturdy.com/api/pkg/organization/service"
	"getsturdy.com/api/pkg/version"

	"github.com/google/uuid"
)

var ErrInvalidLicense = errors.New("invalid license")

type Service struct {
	repo                db.Repository
	organizationService *service_organization.Service
	validator           *validator_license.Validator

	licenseGuard *sync.RWMutex
	// license is the latest license that was retrieved from the licensing server.
	license *licenses.License
}

func New(
	repo db.Repository,
	organizationService *service_organization.Service,
	validator *validator_license.Validator,
) *Service {
	return &Service{
		repo:                repo,
		organizationService: organizationService,
		validator:           validator,

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

func (svc *Service) UpdateLicenseKey(ctx context.Context, key string) error {
	ins, err := svc.Get(ctx)
	if err != nil {
		return fmt.Errorf("could not get installation: %w", err)
	}

	ins.LicenseKey = &key

	// re-load the key immediately
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	license, err := svc.validator.Validate(ctx, *ins.LicenseKey)
	if err != nil {
		return fmt.Errorf("failed to validate license: %w", err)
	}
	// if license.Status != licenses.StatusValid {
	// 	return ErrInvalidLicense
	// }

	if err := svc.UpdateLicense(ctx, license); err != nil {
		return fmt.Errorf("failed to update license: %w", err)
	}

	// save the license in the db
	if err := svc.repo.Update(ctx, ins); err != nil {
		return fmt.Errorf("could not update license key: %w", err)
	}

	return nil
}
