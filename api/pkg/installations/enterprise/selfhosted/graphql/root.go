package graphql

import (
	"context"
	"errors"

	gqlerror "getsturdy.com/api/pkg/graphql/errors"
	"getsturdy.com/api/pkg/graphql/resolvers"
	"getsturdy.com/api/pkg/installations/oss/graphql"
	service_installations "getsturdy.com/api/pkg/installations/service"
)

type rootResolver struct {
	service         *service_installations.Service
	ossRootResolver resolvers.InstallationsRootResolver
}

func New(
	service *service_installations.Service,
	licenseResolver resolvers.LicenseRootResolver,
) resolvers.InstallationsRootResolver {
	return &rootResolver{
		service:         service,
		ossRootResolver: graphql.New(service, licenseResolver),
	}
}

func (r *rootResolver) Installation(ctx context.Context) (resolvers.InstallationsResolver, error) {
	return r.ossRootResolver.Installation(ctx)
}

func (r *rootResolver) UpdateInstallation(ctx context.Context, args resolvers.UpdateInstallationArgs) (resolvers.InstallationsResolver, error) {
	if args.Input.LicenseKey != nil {
		err := r.service.UpdateLicenseKey(ctx, *args.Input.LicenseKey)
		switch {
		case err == nil:
		case errors.Is(err, service_installations.ErrInvalidLicense):
			return nil, gqlerror.Error(err, "message", "The license key is invalid, and can not be used")
		default:
			return nil, gqlerror.Error(err)
		}
	}
	return r.ossRootResolver.Installation(ctx)
}
