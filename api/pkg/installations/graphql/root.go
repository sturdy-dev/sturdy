package graphql

import (
	"getsturdy.com/api/pkg/graphql/resolvers"
	"getsturdy.com/api/pkg/installations"
	service_installations "getsturdy.com/api/pkg/installations/service"
)

type rootResolver struct {
	installation *installations.Installation
	service      *service_installations.Service
}

func New(installation *installations.Installation) resolvers.ServerStatusRootResolver {
	return &rootResolver{
		installation: installation,
	}
}

func (r *rootResolver) ServerStatus() (resolvers.ServerStatusResolver, error) {
	return &resolver{root: r}, nil
}
