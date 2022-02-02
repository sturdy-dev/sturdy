package graphql

import (
	"context"

	gqlerrors "getsturdy.com/api/pkg/graphql/errors"
	"getsturdy.com/api/pkg/graphql/resolvers"
)

type rootResolver struct{}

func New() resolvers.InstallationsRootResolver {
	return &rootResolver{}
}

func (r *rootResolver) Installation(context.Context) (resolvers.InstallationsResolver, error) {
	return nil, gqlerrors.ErrNotImplemented
}

func (r *rootResolver) UpdateInstallation(_ context.Context, _ resolvers.UpdateInstallationArgs) (resolvers.InstallationsResolver, error) {
	return nil, gqlerrors.ErrNotImplemented
}
