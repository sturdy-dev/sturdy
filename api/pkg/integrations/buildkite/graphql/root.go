package graphql

import (
	"context"

	gqlerrors "getsturdy.com/api/pkg/graphql/errors"
	"getsturdy.com/api/pkg/graphql/resolvers"
)

type rootResolver struct{}

func New() resolvers.BuildkiteInstantIntegrationRootResolver {
	return &rootResolver{}
}

func (root *rootResolver) CreateOrUpdateBuildkiteIntegration(ctx context.Context, args resolvers.CreateOrUpdateBuildkiteIntegrationArgs) (resolvers.IntegrationResolver, error) {
	return nil, gqlerrors.ErrNotImplemented
}

func (r *rootResolver) InternalBuildkiteConfigurationByIntegrationID(ctx context.Context, integrationID string) (resolvers.BuildkiteConfigurationResolver, error) {
	return nil, gqlerrors.ErrNotImplemented
}
