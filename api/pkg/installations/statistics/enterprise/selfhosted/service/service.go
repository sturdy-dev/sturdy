package service

import (
	"context"
	"fmt"
	"time"

	service_codebases "getsturdy.com/api/pkg/codebase/service"
	"getsturdy.com/api/pkg/installations"
	"getsturdy.com/api/pkg/installations/statistics"
	service_users "getsturdy.com/api/pkg/user/service"
)

type Service struct {
	installation *installations.Installation

	codebasesService *service_codebases.Service
	usersService     service_users.Service
}

func New(
	installation *installations.Installation,

	codebasesService *service_codebases.Service,
	usersService service_users.Service,
) *Service {
	return &Service{
		installation:     installation,
		codebasesService: codebasesService,
		usersService:     usersService,
	}
}

func (s *Service) Get(ctx context.Context) (*statistics.Statistic, error) {
	usersCount, err := s.usersService.UsersCount(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get users count: %w", err)
	}

	codebasesCount, err := s.codebasesService.CodebaseCount(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get codebases count: %w", err)
	}

	return &statistics.Statistic{
		InstallationID: s.installation.ID,
		LicenseKey:     s.installation.LicenseKey,
		Version:        s.installation.Version,
		RecordedAt:     time.Now(),
		UsersCount:     usersCount,
		CodebasesCount: codebasesCount,
	}, nil
}
