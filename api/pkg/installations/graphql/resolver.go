package graphql

import (
	"context"
	"database/sql"
	"errors"

	gqlerrors "getsturdy.com/api/pkg/graphql/errors"
	"getsturdy.com/api/pkg/graphql/resolvers"
	"getsturdy.com/api/pkg/installations"

	"github.com/graph-gophers/graphql-go"
)

type resolver struct {
	root         *RootResolver
	installation *installations.Installation
}

func (r *resolver) ID() graphql.ID {
	return graphql.ID(r.installation.ID)
}

func (r *resolver) NeedsFirstTimeSetup(ctx context.Context) (bool, error) {
	_, err := r.root.organizationService.GetFirst(ctx)
	switch {
	case err == nil:
		return true, nil
	case errors.Is(err, sql.ErrNoRows):
		return false, nil
	default:
		return false, gqlerrors.Error(err)
	}
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
