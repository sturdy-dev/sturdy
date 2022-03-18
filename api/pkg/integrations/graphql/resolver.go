package graphql

import (
	"context"
	"fmt"

	"getsturdy.com/api/pkg/integrations"
	"getsturdy.com/api/pkg/integrations/providers"

	"getsturdy.com/api/pkg/graphql/resolvers"

	"github.com/graph-gophers/graphql-go"
)

type instantIntegrationProvider struct {
	root        *rootResolver
	integration *integrations.Integration
}

func (ir *instantIntegrationProvider) ID() graphql.ID {
	return graphql.ID(ir.integration.ID)
}

func (ir *instantIntegrationProvider) CodebaseID() graphql.ID {
	return graphql.ID(ir.integration.CodebaseID)
}

func (ir *instantIntegrationProvider) CreatedAt() int32 {
	return int32(ir.integration.CreatedAt.Unix())
}

func (ir *instantIntegrationProvider) UpdatedAt() *int32 {
	if ir.integration.UpdatedAt.IsZero() {
		return nil
	}
	ts := int32(ir.integration.UpdatedAt.Unix())
	return &ts
}

func (ir *instantIntegrationProvider) DeletedAt() *int32 {
	if ir.integration.DeletedAt == nil || ir.integration.DeletedAt.IsZero() {
		return nil
	}
	ts := int32(ir.integration.DeletedAt.Unix())
	return &ts
}

func (ir *instantIntegrationProvider) Provider() (resolvers.InstantIntegrationProviderType, error) {
	switch ir.integration.Provider {
	case providers.ProviderNameBuildkite:
		return resolvers.InstantIntegrationProviderBuildkite, nil
	default:
		return resolvers.InstantIntegrationProviderUndefined, fmt.Errorf("invalid provider: %s", ir.integration.Provider)
	}
}

func (ir *instantIntegrationProvider) ToBuildkiteIntegration() (resolvers.BuildkiteIntegration, bool) {
	if ir.integration.Provider != providers.ProviderNameBuildkite {
		return nil, false
	}
	return &buildkiteProviderResolver{ir}, true
}

type buildkiteProviderResolver struct {
	*instantIntegrationProvider
}

func (br *buildkiteProviderResolver) Configuration(ctx context.Context) (resolvers.BuildkiteConfigurationResolver, error) {
	return br.root.buildkiteRootResolver.InternalBuildkiteConfigurationByIntegrationID(ctx, br.integration.ID)
}
