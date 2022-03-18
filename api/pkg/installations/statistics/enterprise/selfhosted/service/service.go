package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	service_codebases "getsturdy.com/api/pkg/codebases/service"
	service_installations "getsturdy.com/api/pkg/installations/service"
	"getsturdy.com/api/pkg/installations/statistics"
	"getsturdy.com/api/pkg/installations/statistics/enterprise/selfhosted/publisher"
	service_users "getsturdy.com/api/pkg/users/service"
)

type Service struct {
	installationsService *service_installations.Service
	codebasesService     *service_codebases.Service
	usersService         service_users.Service
	publisher            *publisher.Publisher
}

func New(
	installationsService *service_installations.Service,
	codebasesService *service_codebases.Service,
	usersService service_users.Service,

	publisher *publisher.Publisher,
) *Service {
	return &Service{
		installationsService: installationsService,
		codebasesService:     codebasesService,
		usersService:         usersService,
		publisher:            publisher,
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

	ins, err := s.installationsService.Get(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get the current installation: %w", err)
	}

	statistic := &statistics.Statistic{
		InstallationID: ins.ID,
		LicenseKey:     ins.LicenseKey,
		Version:        ins.Version,
		RecordedAt:     time.Now(),
		UsersCount:     usersCount,
		CodebasesCount: codebasesCount,
	}

	if firstUser, err := s.usersService.GetFirstUser(ctx); errors.Is(err, service_users.ErrNotFound) {
		// do nothing
	} else if err != nil {
		return nil, fmt.Errorf("failed to get first user: %w", err)
	} else {
		statistic.FirstUserEmail = &firstUser.Email
	}

	return statistic, nil
}

func (svc *Service) Publish(ctx context.Context) error {
	stats, err := svc.Get(ctx)
	if err != nil {
		return fmt.Errorf("failed to get statistics: %w", err)
	}
	if err := svc.publisher.Publish(ctx, stats); err != nil {
		return fmt.Errorf("failed to publish statistics: %w", err)
	}
	return nil
}
