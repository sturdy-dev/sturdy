package graphql

import (
	"context"

	gqlerror "getsturdy.com/api/pkg/graphql/errors"
	"getsturdy.com/api/pkg/graphql/resolvers"
	"getsturdy.com/api/pkg/installations/enterprise/selfhosted/service"
	"getsturdy.com/api/pkg/installations/graphql"
)

type rootResolver struct {
	*graphql.RootResolver

	service *service.Service
}

func New(
	service *service.Service,
	licenseResolver resolvers.LicenseRootResolver,
	ossResolver *graphql.RootResolver,
) *rootResolver {
	return &rootResolver{
		RootResolver: ossResolver,
		service:      service,
	}
}

func (r *rootResolver) UpdateInstallation(ctx context.Context, args resolvers.UpdateInstallationArgs) (resolvers.InstallationsResolver, error) {
	if args.Input.LicenseKey != nil {
		if err := r.service.UpdateLicenseKey(ctx, *args.Input.LicenseKey); err != nil {
			return nil, gqlerror.Error(err)
		}
	}
	return r.Installation(ctx)
}
