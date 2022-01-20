package graphql

import (
	"context"

	gqlerrors "getsturdy.com/api/pkg/graphql/errors"
	"getsturdy.com/api/pkg/graphql/resolvers"
)

type resolver struct{}

func New() resolvers.ServerStatusRootResolver {
	return resolver{}
}

func (resolver) ServerStatus(_ context.Context) (resolvers.ServerStatusResolver, error) {
	return nil, gqlerrors.ErrNotImplemented
}
