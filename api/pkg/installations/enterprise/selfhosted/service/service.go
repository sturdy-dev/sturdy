package service

import (
	"context"
	"fmt"
	"time"

	"getsturdy.com/api/pkg/installations/db"
	"getsturdy.com/api/pkg/installations/service"
	service_statistics "getsturdy.com/api/pkg/installations/statistics/enterprise/selfhosted/service"
	"getsturdy.com/api/pkg/licenses"
	"getsturdy.com/api/pkg/licenses/enterprise/selfhosted/validator"
)

type Service struct {
	*service.Service

	repo      db.Repository
	validator *validator.Validator

	statisticsService *service_statistics.Service
}

func New(
	s *service.Service,
	valivalidator *validator.Validator,
	repo db.Repository,
	statisticsService *service_statistics.Service,
) *Service {
	return &Service{
		Service:           s,
		validator:         valivalidator,
		repo:              repo,
		statisticsService: statisticsService,
	}
}

func (svc *Service) UpdateLicenseKey(ctx context.Context, key string) error {
	ins, err := svc.Get(ctx)
	if err != nil {
		return fmt.Errorf("could not get installation: %w", err)
	}

	ins.LicenseKey = &key

	// re-load the key immediately
	if _, err := svc.refresh(ctx, key); err != nil {
		return fmt.Errorf("failed to validate license: %w", err)
	}

	if err := svc.repo.Update(ctx, ins); err != nil {
		return fmt.Errorf("could not update license key: %w", err)
	}

	// Send statistics right away
	if err := svc.statisticsService.Publish(ctx); err != nil {
		return fmt.Errorf("could not send installation statistics: %w", err)
	}

	// refresh license status from server again (after published stats)
	if _, err := svc.refresh(ctx, key); err != nil {
		return fmt.Errorf("failed to validate license: %w", err)
	}

	return nil
}

func (svc *Service) refresh(ctx context.Context, licenseKey string) (*licenses.License, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	license, err := svc.validator.Validate(ctx, licenseKey)
	if err != nil {
		return nil, fmt.Errorf("failed to validate license: %w", err)
	}

	if err := svc.UpdateLicense(ctx, license); err != nil {
		return nil, fmt.Errorf("failed to update license: %w", err)
	}

	return license, nil
}
