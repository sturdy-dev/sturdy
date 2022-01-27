package graphql

import (
	"context"

	gqlerrors "getsturdy.com/api/pkg/graphql/errors"
	"getsturdy.com/api/pkg/graphql/resolvers"
	"getsturdy.com/api/pkg/installations"

	"github.com/graph-gophers/graphql-go"
)

type resolver struct {
	root         *rootResolver
	installation *installations.Installation
}

func (r *resolver) ID() graphql.ID {
	return graphql.ID(r.installation.ID)
}

func (r *resolver) NeedsFirstTimeSetup(ctx context.Context) (bool, error) {
	hasOrg, err := r.root.service.HasOrganization(ctx)
	if err != nil {
		return false, gqlerrors.Error(err)
	}
	return hasOrg, nil
}

func (r *resolver) Version() string {
	return r.installation.Version
}

func (r *resolver) License(ctx context.Context) (resolvers.LicenseResolver, error) {
	if r.installation.LicenseKey == nil {
		return nil, gqlerrors.ErrNotFound
	}
	return r.root.licenseResolver.InternalByKey(ctx, *r.installation.LicenseKey)
}
