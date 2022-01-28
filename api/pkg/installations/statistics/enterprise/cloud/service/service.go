package service

import (
	"context"
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
