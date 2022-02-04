package service

import (
	"context"
	"fmt"
	"time"

	"getsturdy.com/api/pkg/installations/db"
	"getsturdy.com/api/pkg/installations/service"
	"getsturdy.com/api/pkg/licenses/enterprise/selfhosted/validator"
)

type Service struct {
	*service.Service

	repo      db.Repository
	validator *validator.Validator
}

func New(
	s *service.Service,
	valivalidator *validator.Validator,
	repo db.Repository,
) *Service {
	return &Service{
		Service:   s,
		validator: valivalidator,
		repo:      repo,
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

	if err := svc.UpdateLicense(ctx, license); err != nil {
		return fmt.Errorf("failed to update license: %w", err)
	}

	if err := svc.repo.Update(ctx, ins); err != nil {
		return fmt.Errorf("could not update license key: %w", err)
	}

	return nil
}
