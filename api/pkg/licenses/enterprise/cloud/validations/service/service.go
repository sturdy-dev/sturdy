package service

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"

	"getsturdy.com/api/pkg/licenses"
	"getsturdy.com/api/pkg/licenses/enterprise/cloud/validations"
	"getsturdy.com/api/pkg/licenses/enterprise/cloud/validations/db"
)

type Service struct {
	repo db.Repository
}

func New(
	repo db.Repository,
) *Service {
	return &Service{
		repo: repo,
	}
}

func (svc *Service) Create(ctx context.Context, licenseID licenses.ID, status licenses.Status) (*validations.Validation, error) {
	validation := &validations.Validation{
		ID:        uuid.New().String(),
		LicenseID: licenseID,
		Timestamp: time.Now(),
		Status:    status,
	}

	if err := svc.repo.Create(ctx, validation); err != nil {
		return nil, fmt.Errorf("failed to record validation: %w", err)
	}

	return validation, nil
}
