package graphql

import (
	"context"
	"fmt"

	service_auth "getsturdy.com/api/pkg/auth/service"
	"getsturdy.com/api/pkg/changes"
	service_change "getsturdy.com/api/pkg/changes/service"
	"getsturdy.com/api/pkg/ci/service"
	"getsturdy.com/api/pkg/codebases"
	gqlerrors "getsturdy.com/api/pkg/graphql/errors"
	"getsturdy.com/api/pkg/graphql/resolvers"
	"getsturdy.com/api/pkg/integrations"
	"getsturdy.com/api/pkg/integrations/providers"
	service_workspaces "getsturdy.com/api/pkg/workspaces/service"
)

type rootResolver struct {
	svc              *service.Service
	changeService    *service_change.Service
	authService      *service_auth.Service
	workspaceService *service_workspaces.Service

	buildkiteRootResolver resolvers.BuildkiteInstantIntegrationRootResolver
	statusesRootResolver  resolvers.StatusesRootResolver
}

func NewRootResolver(
	svc *service.Service,
	changeService *service_change.Service,
	authService *service_auth.Service,
	workspaceService *service_workspaces.Service,

	buildkiteRootResolver resolvers.BuildkiteInstantIntegrationRootResolver,
	statusesRootResolver resolvers.StatusesRootResolver,
) resolvers.IntegrationRootResolver {
	return &rootResolver{
		svc:              svc,
		changeService:    changeService,
		authService:      authService,
		workspaceService: workspaceService,

		buildkiteRootResolver: buildkiteRootResolver,
		statusesRootResolver:  statusesRootResolver,
	}
}

func (r *rootResolver) TriggerInstantIntegration(ctx context.Context, args resolvers.TriggerInstantIntegrationArgs) ([]resolvers.StatusResolver, error) {
	if args.Input.ChangeID != nil {
		return r.triggerInstantIntegrationChangeID(ctx, changes.ID(*args.Input.ChangeID), args.Input.Providers)
	} else if args.Input.WorkspaceID != nil {
		return r.triggerInstantIntegrationWorkspaceID(ctx, string(*args.Input.WorkspaceID), args.Input.Providers)
	} else {
		return nil, gqlerrors.Error(gqlerrors.ErrBadRequest, "one of change id or workspace id must be set")
	}
}

func (r *rootResolver) triggerInstantIntegrationWorkspaceID(ctx context.Context, workspaceID string, providers *[]resolvers.InstantIntegrationProviderType) ([]resolvers.StatusResolver, error) {
	ws, err := r.workspaceService.GetByID(ctx, workspaceID)
	if err != nil {
		return nil, gqlerrors.Error(err)
	}

	if err := r.authService.CanWrite(ctx, ws); err != nil {
		return nil, gqlerrors.Error(err)
	}

	var triggerOptions []service.TriggerOption
	if providers != nil {
		for _, provider := range *providers {
			providerType, err := convertProvider(provider)
			if err != nil {
				return nil, gqlerrors.Error(gqlerrors.ErrBadRequest, "integrations", err.Error())
			}
			triggerOptions = append(triggerOptions, service.WithProvider(providerType))
		}
	}

	ss, err := r.svc.TriggerWorkspace(ctx, ws, triggerOptions...)
	if err != nil {
		return nil, gqlerrors.Error(err)
	}

	rr := make([]resolvers.StatusResolver, 0, len(ss))
	for _, s := range ss {
		rr = append(rr, r.statusesRootResolver.InternalStatus(s))
	}

	return rr, nil
}

func (r *rootResolver) triggerInstantIntegrationChangeID(ctx context.Context, changeID changes.ID, providers *[]resolvers.InstantIntegrationProviderType) ([]resolvers.StatusResolver, error) {
	ch, err := r.changeService.GetChangeByID(ctx, changeID)
	if err != nil {
		return nil, gqlerrors.Error(err)
	}

	if err := r.authService.CanWrite(ctx, ch); err != nil {
		return nil, gqlerrors.Error(err)
	}

	var triggerOptions []service.TriggerOption
	if providers != nil {
		for _, provider := range *providers {
			providerName, err := convertProvider(provider)
			if err != nil {
				return nil, gqlerrors.Error(gqlerrors.ErrBadRequest, "integrations", err.Error())
			}
			triggerOptions = append(triggerOptions, service.WithProvider(providerName))
		}
	}

	ss, err := r.svc.TriggerChange(ctx, ch, triggerOptions...)
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

	if err := r.authService.CanWrite(ctx, &codebases.Codebase{ID: cfg.CodebaseID}); err != nil {
		return nil, gqlerrors.Error(err)
	}

	if err := r.svc.Delete(ctx, string(args.Input.ID)); err != nil {
		return nil, gqlerrors.Error(err)
	}

	return r.InternalIntegrationByID(ctx, string(args.Input.ID))
}

func convertProvider(in resolvers.InstantIntegrationProviderType) (providers.ProviderName, error) {
	switch in {
	case resolvers.InstantIntegrationProviderBuildkite:
		return providers.ProviderNameBuildkite, nil
	default:
		return providers.ProviderNameUndefined, fmt.Errorf("invalid provider: %s", in)
	}
}

func (r *rootResolver) InternalIntegrationsByCodebaseID(ctx context.Context, codebaseID codebases.ID) ([]resolvers.IntegrationResolver, error) {
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
