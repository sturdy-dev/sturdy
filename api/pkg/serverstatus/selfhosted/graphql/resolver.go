package graphql

import (
	"context"

	"github.com/graph-gophers/graphql-go"

	service_serverstatus "getsturdy.com/api/pkg/serverstatus/selfhosted/service"

	gqlerrors "getsturdy.com/api/pkg/graphql/errors"
	"getsturdy.com/api/pkg/graphql/resolvers"
	"getsturdy.com/api/pkg/version"
)

type rootResolver struct {
	service *service_serverstatus.Service
}

func New(service *service_serverstatus.Service) resolvers.ServerStatusRootResolver {
	return &rootResolver{
		service: service,
	}
}

func (r *rootResolver) ServerStatus(ctx context.Context) (resolvers.ServerStatusResolver, error) {
	return &resolver{root: r}, nil
}

type resolver struct {
	root *rootResolver
}

func (resolver) ID() graphql.ID {
	return "_server"
}

func (r *resolver) NeedsFirstTimeSetup(ctx context.Context) (bool, error) {
	hasOrg, err := r.root.service.HasOrganization(ctx)
	if err != nil {
		return false, gqlerrors.Error(err)
	}
	if hasOrg {
		return false, nil
	}
	return true, nil
}

func (resolver) Version() string {
	return version.Version
}
