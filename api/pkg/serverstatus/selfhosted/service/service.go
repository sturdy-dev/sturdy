package service

import (
	"context"
	"database/sql"

	"github.com/pkg/errors"

	service_organization "getsturdy.com/api/pkg/organization/service"
)

type Service struct {
	organizationService *service_organization.Service
}

func New(organizationService *service_organization.Service) *Service {
	return &Service{
		organizationService: organizationService,
	}
}

func (svc *Service) HasOrganization(ctx context.Context) (bool, error) {
	_, err := svc.organizationService.GetFirst(ctx)
	switch {
	case err == nil:
		return true, nil
	case errors.Is(err, sql.ErrNoRows):
		return false, nil
	default:
		return false, err
	}
}
