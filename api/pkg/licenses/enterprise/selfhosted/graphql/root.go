package graphql

import (
	"context"
	"fmt"

	gqlerrors "getsturdy.com/api/pkg/graphql/errors"
	"getsturdy.com/api/pkg/graphql/resolvers"
	"getsturdy.com/api/pkg/installations"
	"getsturdy.com/api/pkg/licenses/enterprise/selfhosted/validator"
)

type licenseRootResovler struct {
	validator    *validator.Validator
	installation installations.GetInstallationFunc
}

func New(
	validator *validator.Validator,
	installation installations.GetInstallationFunc,
) resolvers.LicenseRootResolver {
	return &licenseRootResovler{
		validator:    validator,
		installation: installation,
	}
}

func (r *licenseRootResovler) InternalByKey(ctx context.Context, key string) (resolvers.LicenseResolver, error) {
	ins, err := r.installation()
	if err != nil {
		return nil, gqlerrors.Error(err)
	}

	if ins.LicenseKey == nil {
		return nil, gqlerrors.ErrNotFound
	}

	license, err := r.validator.Validate(ctx, *ins.LicenseKey)
	if err != nil {
		return nil, gqlerrors.Error(fmt.Errorf("failed to validate license: %w", err))
	}

	return &licenseResolver{
		license: license,
	}, nil
}

func (r *licenseRootResovler) InternalListForOrganizationID(ctx context.Context, organizationID string) ([]resolvers.LicenseResolver, error) {
	ins, err := r.installation()
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
