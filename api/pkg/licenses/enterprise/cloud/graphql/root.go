package graphql

import (
	"context"

	gqlerrors "getsturdy.com/api/pkg/graphql/errors"
	"getsturdy.com/api/pkg/graphql/resolvers"
	service_licenses "getsturdy.com/api/pkg/licenses/enterprise/cloud/service"
)

type licenseRootResovler struct {
	svc *service_licenses.Service
}

func New(
	svc *service_licenses.Service,
) resolvers.LicenseRootResolver {
	return &licenseRootResovler{
		svc: svc,
	}
}

func (r *licenseRootResovler) InternalByKey(ctx context.Context, key string) (resolvers.LicenseResolver, error) {
	l, err := r.svc.ValidateByKey(ctx, key)
	if err != nil {
		return nil, gqlerrors.Error(err)
	}
	return &licenseResolver{
		root:    r,
		license: l,
	}, nil
}

func (r *licenseRootResovler) InternalListForOrganizationID(ctx context.Context, id string) ([]resolvers.LicenseResolver, error) {
	ll, err := r.svc.ListByOrganizationID(ctx, id)
	if err != nil {
		return nil, gqlerrors.Error(err)
	}

	resolvers := make([]resolvers.LicenseResolver, 0, len(ll))
	for _, l := range ll {
		resolvers = append(resolvers, &licenseResolver{
			root:    r,
			license: l,
		})
	}
	return resolvers, nil
}
