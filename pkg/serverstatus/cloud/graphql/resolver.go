package graphql

import (
	"context"

	gqlerrors "mash/pkg/graphql/errors"
	"mash/pkg/graphql/resolvers"
)

type resolver struct{}

func New() resolvers.ServerStatusRootResolver {
	return resolver{}
}

func (resolver) ServerStatus(_ context.Context) (resolvers.ServerStatusResolver, error) {
	return nil, gqlerrors.ErrNotImplemented
}
