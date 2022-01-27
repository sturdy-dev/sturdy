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

func New(installation *installations.Installation) resolvers.InstallationsRootResolver {
	return &rootResolver{
		installation: installation,
	}
}

func (r *rootResolver) Installation() (resolvers.InstallationsResolver, error) {
	return &resolver{root: r}, nil
}
