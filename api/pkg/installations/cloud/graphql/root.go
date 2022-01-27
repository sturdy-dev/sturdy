package graphql

import (
	gqlerrors "getsturdy.com/api/pkg/graphql/errors"
	"getsturdy.com/api/pkg/graphql/resolvers"
)

type rootResolver struct{}

func New() resolvers.InstallationsRootResolver {
	return &rootResolver{}
}

func (r *rootResolver) Installation() (resolvers.InstallationsResolver, error) {
	return nil, gqlerrors.ErrNotImplemented
}
