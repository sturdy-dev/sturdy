package graphql

import (
	"context"

	gqlerrors "getsturdy.com/api/pkg/graphql/errors"
	"getsturdy.com/api/pkg/graphql/resolvers"
	service_installations "getsturdy.com/api/pkg/installations/service"
)

type rootResolver struct {
	service         *service_installations.Service
	licenseResolver resolvers.LicenseRootResolver
}

func New(
	service *service_installations.Service,
	licenseResolver resolvers.LicenseRootResolver,
) resolvers.InstallationsRootResolver {
	return &rootResolver{
		service:         service,
		licenseResolver: licenseResolver,
	}
}

func (r *rootResolver) Installation(ctx context.Context) (resolvers.InstallationsResolver, error) {
	installation, err := r.service.Get(ctx)
	if err != nil {
		return nil, gqlerrors.Error(err)
	}
	return &resolver{
		root:         r,
		installation: installation,
	}, nil
}

func (r *rootResolver) UpdateInstallation(_ context.Context, _ resolvers.UpdateInstallationArgs) (resolvers.InstallationsResolver, error) {
	return nil, gqlerrors.ErrNotImplemented
}
