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

func (r *licenseRootResovler) InternalByKey(context.Context, string) (resolvers.LicenseResolver, error) {
	return nil, gqlerrors.ErrNotImplemented
}

func (*licenseRootResovler) InternalListForOrganizationID(context.Context, string) ([]resolvers.LicenseResolver, error) {
	return nil, gqlerrors.ErrNotImplemented
}
