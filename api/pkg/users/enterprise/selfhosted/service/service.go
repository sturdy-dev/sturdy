package service

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"getsturdy.com/api/pkg/installations"
	service_organization "getsturdy.com/api/pkg/organization/service"
	"getsturdy.com/api/pkg/users"
	"getsturdy.com/api/pkg/users/service"
)

const maxUsersWithoutLicense = 10

type Service struct {
	*service.UserSerice

	organizationService *service_organization.Service
	installation        installations.GetInstallationFunc
}

func New(
	userService *service.UserSerice,
	organizationService *service_organization.Service,
	installation installations.GetInstallationFunc,
) *Service {
	return &Service{
		UserSerice:          userService,
		organizationService: organizationService,
		installation:        installation,
	}
}

func (s *Service) validate(ctx context.Context) error {
	ins, err := s.installation()
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

func (s *Service) Create(ctx context.Context, name, email string) (*users.User, error) {
	if err := s.validate(ctx); err != nil {
		return nil, err
	}
	return s.UserSerice.Create(ctx, name, email)
}

func (s *Service) CreateWithPassword(ctx context.Context, name, password, email string) (*users.User, error) {
	if err := s.validate(ctx); err != nil {
		return nil, err
	}

	usr, err := s.UserSerice.CreateWithPassword(ctx, name, password, email)
	if err != nil {
		return nil, err
	}

	// If this instance has an organization, auto-add this user
	first, err := s.organizationService.GetFirst(ctx)
	switch {
	case err == nil:
		// add this user
		if _, err := s.organizationService.AddMember(ctx, first.ID, usr.ID, usr.ID); err != nil {
			return nil, fmt.Errorf("failed to add member to existing org: %w", err)
		}
	case errors.Is(err, sql.ErrNoRows):
	// first org has not been created yet, this user will create it later
	case err != nil:
		return nil, fmt.Errorf("failed to check if an organization already exists: %w", err)
	}

	return usr, nil
}
