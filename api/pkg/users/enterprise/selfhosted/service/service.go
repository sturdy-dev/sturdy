package service

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/pkg/errors"

	service_organization "getsturdy.com/api/pkg/organization/service"
	"getsturdy.com/api/pkg/users"
	"getsturdy.com/api/pkg/users/service"
)

type Service struct {
	*service.UserSerice
	organizationService *service_organization.Service
}

func New(
	userService *service.UserSerice,
	organizationService *service_organization.Service,
) *Service {
	return &Service{
		UserSerice:          userService,
		organizationService: organizationService,
	}
}

func (s *Service) CreateWithPassword(ctx context.Context, name, password, email string) (*user.User, error) {
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
