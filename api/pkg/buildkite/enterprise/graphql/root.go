package graphql

import (
	"context"
	"fmt"
	"time"

	service_auth "getsturdy.com/api/pkg/auth/service"
	"getsturdy.com/api/pkg/buildkite"
	service_buildkite "getsturdy.com/api/pkg/buildkite/enterprise/service"
	service_ci "getsturdy.com/api/pkg/ci/service"
	"getsturdy.com/api/pkg/codebases"
	gqlerrors "getsturdy.com/api/pkg/graphql/errors"
	"getsturdy.com/api/pkg/graphql/resolvers"
	"getsturdy.com/api/pkg/integrations"
	"getsturdy.com/api/pkg/integrations/providers"

	"github.com/google/uuid"
)

type rootResolver struct {
	authService                    *service_auth.Service
	buildkiteService               *service_buildkite.Service
	instantIntegrationService      *service_ci.Service
	instantIntegrationRootResolver *resolvers.IntegrationRootResolver
}

var seedFiles = []string{
	"buildkite.yml",
	"buildkite.yaml",
	"buildkite.json",
	".buildkite/pipeline.yml",
	".buildkite/pipeline.yaml",
	".buildkite/pipeline.json",
	"buildkite/pipeline.yml",
	"buildkite/pipeline.yaml",
	"buildkite/pipeline.json",
}

func New(
	authService *service_auth.Service,
	buildkiteService *service_buildkite.Service,
	instantIntegrationService *service_ci.Service,
	instantIntegrationRootResolver *resolvers.IntegrationRootResolver,
) resolvers.BuildkiteInstantIntegrationRootResolver {
	return &rootResolver{
		authService:                    authService,
		buildkiteService:               buildkiteService,
		instantIntegrationService:      instantIntegrationService,
		instantIntegrationRootResolver: instantIntegrationRootResolver,
	}
}

func (root *rootResolver) createNewConfiguration(ctx context.Context, args resolvers.CreateOrUpdateBuildkiteIntegrationArgs) (*integrations.Integration, error) {
	integration := &integrations.Integration{
		ID:           uuid.NewString(),
		CodebaseID:   codebases.ID(args.Input.CodebaseID),
		Provider:     providers.ProviderNameBuildkite,
		ProviderType: providers.ProviderTypeBuild,
		CreatedAt:    time.Now(),
		SeedFiles:    seedFiles,
	}

	if err := root.instantIntegrationService.CreateIntegration(ctx, integration); err != nil {
		return nil, fmt.Errorf("failed to create integration: %w", err)
	}

	cfg := &buildkite.Config{
		ID:               uuid.NewString(),
		IntegrationID:    integration.ID,
		CodebaseID:       codebases.ID(args.Input.CodebaseID),
		OrganizationName: args.Input.OrganizationName,
		PipelineName:     args.Input.PipelineName,
		APIToken:         args.Input.APIToken,
		WebhookSecret:    args.Input.WebhookSecret,
		CreatedAt:        time.Now(),
	}

	if err := root.buildkiteService.CreateIntegration(ctx, cfg); err != nil {
		return nil, fmt.Errorf("failed to create configuration: %w", err)
	}

	return integration, nil
}

func (root *rootResolver) updateConfiguration(ctx context.Context, existingCfg *buildkite.Config, args resolvers.CreateOrUpdateBuildkiteIntegrationArgs) (*integrations.Integration, error) {
	existingIntegrations, err := root.instantIntegrationService.ListByCodebaseID(ctx, codebases.ID(args.Input.CodebaseID))
	if err != nil {
		return nil, fmt.Errorf("failed to list existing integrations: %w", err)
	}

	var integration *integrations.Integration
	for _, i := range existingIntegrations {
		if i.Provider == providers.ProviderNameBuildkite {
			integration = i
		}
	}

	if integration == nil {
		return nil, fmt.Errorf("integraion must exist, but it was not found: %w", err)
	}

	configChanged := existingCfg.OrganizationName != args.Input.OrganizationName ||
		existingCfg.PipelineName != args.Input.PipelineName ||
		existingCfg.APIToken != args.Input.APIToken ||
		existingCfg.WebhookSecret != args.Input.WebhookSecret

	if !configChanged {
		return integration, nil
	}

	existingCfg.PipelineName = args.Input.PipelineName
	existingCfg.OrganizationName = args.Input.OrganizationName
	existingCfg.APIToken = args.Input.APIToken
	existingCfg.WebhookSecret = args.Input.WebhookSecret
	existingCfg.UpdatedAt = time.Now()
	if err := root.buildkiteService.UpdateIntegration(ctx, existingCfg); err != nil {
		return nil, fmt.Errorf("failed to update configuration: %w", err)
	}

	integration.UpdatedAt = time.Now()
	if err := root.instantIntegrationService.UpdateIntegration(ctx, integration); err != nil {
		return nil, fmt.Errorf("failed to update integration: %w", err)
	}

	return integration, nil
}

func (root *rootResolver) CreateOrUpdateBuildkiteIntegration(ctx context.Context, args resolvers.CreateOrUpdateBuildkiteIntegrationArgs) (resolvers.IntegrationResolver, error) {
	if err := root.authService.CanWrite(ctx, &codebases.Codebase{ID: codebases.ID(args.Input.CodebaseID)}); err != nil {
		return nil, gqlerrors.Error(err)
	}

	// Create new
	if args.Input.IntegrationID == nil {
		integration, err := root.createNewConfiguration(ctx, args)
		if err != nil {
			return nil, gqlerrors.Error(fmt.Errorf("failed to create new configuration: %w", err))
		}
		return (*root.instantIntegrationRootResolver).InternalIntegrationProvider(integration), nil
	}

	// Update existing
	existingCfg, err := root.buildkiteService.GetConfigurationByIntegrationID(ctx, string(*args.Input.IntegrationID))
	if err != nil {
		return nil, gqlerrors.Error(err)
	}
	integration, err := root.updateConfiguration(ctx, existingCfg, args)
	if err != nil {
		return nil, gqlerrors.Error(fmt.Errorf("failed to update existing configuration: %w", err))
	}

	return (*root.instantIntegrationRootResolver).InternalIntegrationProvider(integration), nil
}

func (r *rootResolver) InternalBuildkiteConfigurationByIntegrationID(ctx context.Context, integrationID string) (resolvers.BuildkiteConfigurationResolver, error) {
	cfg, err := r.buildkiteService.GetConfigurationByIntegrationID(ctx, integrationID)
	if err != nil {
		return nil, gqlerrors.Error(err)
	}
	return &buildkiteConfigurationResover{
		buildkiteConfig: cfg,
	}, nil
}
