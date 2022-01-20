package graphql

import (
	"context"
	"fmt"
	"getsturdy.com/api/pkg/codebase"
	"getsturdy.com/api/pkg/integrations"

	service_auth "getsturdy.com/api/pkg/auth/service"
	"getsturdy.com/api/pkg/change"
	db_change "getsturdy.com/api/pkg/change/db"
	"getsturdy.com/api/pkg/ci/service"
	gqlerrors "getsturdy.com/api/pkg/graphql/errors"
	"getsturdy.com/api/pkg/graphql/resolvers"
)

type rootResolver struct {
	svc         *service.Service
	changeRepo  db_change.Repository
	authService *service_auth.Service

	buildkiteRootResolver resolvers.BuildkiteInstantIntegrationRootResolver
	statusesRootResolver  resolvers.StatusesRootResolver
}

func NewRootResolver(
	svc *service.Service,
	changeRepo db_change.Repository,
	authService *service_auth.Service,

	buildkiteRootResolver resolvers.BuildkiteInstantIntegrationRootResolver,
	statusesRootResolver resolvers.StatusesRootResolver,
) resolvers.IntegrationRootResolver {
	return &rootResolver{
		svc:         svc,
		changeRepo:  changeRepo,
		authService: authService,

		buildkiteRootResolver: buildkiteRootResolver,
		statusesRootResolver:  statusesRootResolver,
	}
}

func (r *rootResolver) TriggerInstantIntegration(ctx context.Context, args resolvers.TriggerInstantIntegrationArgs) ([]resolvers.StatusResolver, error) {
	ch, err := r.changeRepo.Get(change.ID(args.Input.ChangeID))
	if err != nil {
		return nil, gqlerrors.Error(err)
	}

	if err := r.authService.CanWrite(ctx, ch); err != nil {
		return nil, gqlerrors.Error(err)
	}

	triggerOptions := []service.TriggerOption{}
	if args.Input.Providers != nil {
		for _, provider := range *args.Input.Providers {
			providerType, err := convertProvider(provider)
			if err != nil {
				return nil, gqlerrors.Error(gqlerrors.ErrBadRequest, "integrations", err.Error())
			}
			triggerOptions = append(triggerOptions, service.WithProvider(providerType))
		}
	}

	ss, err := r.svc.Trigger(ctx, &ch, triggerOptions...)
	if err != nil {
		return nil, gqlerrors.Error(err)
	}
	rr := make([]resolvers.StatusResolver, 0, len(ss))
	for _, s := range ss {
		rr = append(rr, r.statusesRootResolver.InternalStatus(s))
	}

	return rr, nil
}

func (r *rootResolver) DeleteIntegration(ctx context.Context, args resolvers.DeleteIntegrationArgs) (resolvers.IntegrationResolver, error) {
	cfg, err := r.svc.GetByID(ctx, string(args.Input.ID))
	if err != nil {
		return nil, gqlerrors.Error(err)
	}

	if err := r.authService.CanWrite(ctx, &codebase.Codebase{ID: cfg.CodebaseID}); err != nil {
		return nil, gqlerrors.Error(err)
	}

	if err := r.svc.Delete(ctx, string(args.Input.ID)); err != nil {
		return nil, gqlerrors.Error(err)
	}

	return r.InternalIntegrationByID(ctx, string(args.Input.ID))
}

func convertProvider(in resolvers.InstantIntegrationProviderType) (integrations.ProviderType, error) {
	switch in {
	case resolvers.InstantIntegrationProviderBuildkite:
		return integrations.ProviderTypeBuildkite, nil
	default:
		return integrations.ProviderTypeUndefined, fmt.Errorf("invalid provider: %s", in)
	}
}

func (r *rootResolver) InternalIntegrationsByCodebaseID(ctx context.Context, codebaseID string) ([]resolvers.IntegrationResolver, error) {
	cfgs, err := r.svc.ListByCodebaseID(ctx, codebaseID)
	if err != nil {
		return nil, gqlerrors.Error(err)
	}

	rr := make([]resolvers.IntegrationResolver, 0, len(cfgs))
	for _, cfg := range cfgs {
		rr = append(rr, &instantIntegrationProvider{
			integration: cfg,
			root:        r,
		})
	}

	return rr, nil
}

func (r *rootResolver) InternalIntegrationByID(ctx context.Context, integrationID string) (resolvers.IntegrationResolver, error) {
	res, err := r.svc.GetByID(ctx, integrationID)
	if err != nil {
		return nil, gqlerrors.Error(err)
	}
	return &instantIntegrationProvider{integration: res, root: r}, nil
}

func (r *rootResolver) InternalIntegrationProvider(integration *integrations.Integration) resolvers.IntegrationResolver {
	return &instantIntegrationProvider{
		root:        r,
		integration: integration,
	}
}
