package graphql

import (
	"context"

	gqlerrors "getsturdy.com/api/pkg/graphql/errors"
	"getsturdy.com/api/pkg/graphql/resolvers"
)

type licenseRootResovler struct{}

func New() resolvers.LicenseRootResolver {
	return &licenseRootResovler{}
}

func (*licenseRootResovler) ValidateLicense(_ context.Context, _ resolvers.ValidateLicenseArgs) (resolvers.LicenseValidation, error) {
	return nil, gqlerrors.ErrNotImplemented
}

func (*licenseRootResovler) InternalListForOrganization(_ context.Context, _ string) ([]resolvers.LicenseResolver, error) {
	return nil, gqlerrors.ErrNotImplemented
}
