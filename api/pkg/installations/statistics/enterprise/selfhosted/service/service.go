package service

import (
	"context"

	"getsturdy.com/api/pkg/installations/statistics"
)

type Service struct{}

func New() *Service {
	return &Service{}
}

func (s *Service) Get(ctx context.Context) (*statistics.Statistic, error) {
	stat := &statistics.Statistic{}
	return stat, nil
}
