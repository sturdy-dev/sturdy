package graphql

import (
	"context"

	"github.com/graph-gophers/graphql-go"

	service_serverstatus "mash/pkg/serverstatus/selfhosted/service"

	gqlerrors "mash/pkg/graphql/errors"
	"mash/pkg/graphql/resolvers"
	"mash/pkg/version"
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
