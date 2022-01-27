package graphql

import (
	gqlerrors "getsturdy.com/api/pkg/graphql/errors"
	"getsturdy.com/api/pkg/graphql/resolvers"
)

type rootResolver struct{}

func New() resolvers.ServerStatusRootResolver {
	return &rootResolver{}
}

func (r *rootResolver) ServerStatus() (resolvers.ServerStatusResolver, error) {
	return nil, gqlerrors.ErrNotImplemented
}
