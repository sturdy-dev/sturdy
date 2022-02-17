package graphql

import (
	"context"
	"fmt"

	gqlerrors "getsturdy.com/api/pkg/graphql/errors"
	"getsturdy.com/api/pkg/graphql/resolvers"
	service_installations "getsturdy.com/api/pkg/installations/service"
	"getsturdy.com/api/pkg/licenses/enterprise/selfhosted/validator"
)

type licenseRootResovler struct {
	validator            *validator.Validator
	installationsService *service_installations.Service
}

func New(
	validator *validator.Validator,
	installationsService *service_installations.Service,
) resolvers.LicenseRootResolver {
	return &licenseRootResovler{
		validator:            validator,
		installationsService: installationsService,
	}
}

func (r *licenseRootResovler) InternalByKey(ctx context.Context, key string) (resolvers.LicenseResolver, error) {
	license, err := r.validator.Validate(ctx, key)
	if err != nil {
		return nil, gqlerrors.Error(fmt.Errorf("failed to validate license: %w", err))
	}

	return &licenseResolver{
		license: license,
	}, nil
}

func (r *licenseRootResovler) InternalListForOrganizationID(ctx context.Context, organizationID string) ([]resolvers.LicenseResolver, error) {
	ins, err := r.installationsService.Get(ctx)
	if err != nil {
		return nil, gqlerrors.Error(err)
	}

	if ins.LicenseKey == nil {
		return nil, nil
	}

	l, err := r.InternalByKey(ctx, *ins.LicenseKey)
	if err != nil {
		return nil, err
	}
	return []resolvers.LicenseResolver{l}, nil
}
