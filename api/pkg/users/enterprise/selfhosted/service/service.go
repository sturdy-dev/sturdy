package service

import (
	"context"
	"fmt"

	service_installations "getsturdy.com/api/pkg/installations/service"
	"getsturdy.com/api/pkg/users"
	service_oss_selfhosted "getsturdy.com/api/pkg/users/oss/selfhosted/service"
	"getsturdy.com/api/pkg/users/service"
)

const maxUsersWithoutLicense = 10

type Service struct {
	*service_oss_selfhosted.Service

	installationService *service_installations.Service
}

func New(
	userService *service_oss_selfhosted.Service,
	installationService *service_installations.Service,
) *Service {
	return &Service{
		Service:             userService,
		installationService: installationService,
	}
}

func (s *Service) ValidateUserCount(ctx context.Context) error {
	ins, err := s.installationService.Get(ctx)
	if err != nil {
		return fmt.Errorf("failed to get current installation: %w", err)
	}
	if ins.License != nil {
		return nil
	}
	usersCount, err := s.UsersCount(ctx)
	if err != nil {
		return fmt.Errorf("failed to get users count: %w", err)
	}
	if usersCount >= maxUsersWithoutLicense {
		return service.ErrExceeded
	}
	return nil
}

func (s *Service) CreateWithPassword(ctx context.Context, name, password, email string) (*users.User, error) {
	if err := s.ValidateUserCount(ctx); err != nil {
		return nil, err
	}

	usr, err := s.Service.CreateWithPassword(ctx, name, password, email)
	if err != nil {
		return nil, err
	}

	return usr, nil
}
