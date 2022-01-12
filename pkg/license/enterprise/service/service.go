package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"

	"mash/pkg/license/enterprise/db"
	"mash/pkg/license/enterprise/license"
)

type Service struct {
	licenseRepository           db.Repository
	licenseValidationRepository db.ValidationRepository
}

func New(
	licenseRepository db.Repository,
	licenseValidationRepository db.ValidationRepository,
) *Service {
	return &Service{
		licenseRepository:           licenseRepository,
		licenseValidationRepository: licenseValidationRepository,
	}
}

// Validate validates the license key, given the user reported data in validation
func (svc *Service) Validate(ctx context.Context, key string, validation license.SelfHostedLicenseValidation) error {
	var status error

	// Make sure that validation attempts of non existing keys are recorded
	if l, err := svc.licenseRepository.Get(ctx, key); err == nil {
		status = validate(l, validation)
	} else {
		status = err
	}

	// Populate validation object with more data
	validation.SelfHostedLicenseID = key // Note: this is user reported data, and this key might not exist
	validation.ID = uuid.NewString()
	validation.ValidatedAt = time.Now()
	validation.Status = status == nil

	// Save validation attempt
	if err := svc.licenseValidationRepository.Record(ctx, validation); err != nil {
		return fmt.Errorf("failed to record validation: %w", err)
	}

	return status
}

var (
	ErrTooManyUsers = errors.New("license validation error too many users")
	ErrExpired      = errors.New("license has expired")
)

func validate(l *license.SelfHostedLicense, validation license.SelfHostedLicenseValidation) error {
	if !l.Active {
		return ErrExpired
	}
	if validation.ReportedUserCount > l.Seats {
		return ErrTooManyUsers
	}
	// TODO: More checks?
	return nil
}

func (svc *Service) ListOrganizationSubscriptions(ctx context.Context, orgID string) ([]*license.SelfHostedLicense, error) {
	return svc.licenseRepository.ListByCloudOrganizationID(ctx, orgID)
}
