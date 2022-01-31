package service

import (
	"context"
	"errors"
	"fmt"

	"getsturdy.com/api/pkg/installations/statistics"
	"getsturdy.com/api/pkg/installations/statistics/enterprise/cloud/db"
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

func (s *Service) Accept(ctx context.Context, statistic *statistics.Statistic) error {
	if err := s.repo.Create(ctx, statistic); err != nil {
		return fmt.Errorf("failed to create statistic: %w", err)
	}
	return nil
}

var (
	ErrNotFound = errors.New("not found")
)

func (s *Service) GetByLicenseKey(ctx context.Context, licenseKey string) (*statistics.Statistic, error) {
	if lisence, err := s.repo.GetByLicenseKey(ctx, licenseKey); errors.Is(err, db.ErrNotFound) {
		return nil, ErrNotFound
	} else if err != nil {
		return nil, fmt.Errorf("failed to get statistic: %w", err)
	} else {
		return lisence, nil
	}
}
