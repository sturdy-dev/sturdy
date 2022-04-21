package grapqhl

import (
	"context"
	"fmt"

	services_auth "getsturdy.com/api/pkg/auth/service"
	gqlerrors "getsturdy.com/api/pkg/graphql/errors"
	"getsturdy.com/api/pkg/graphql/resolvers"
	service_land "getsturdy.com/api/pkg/land/service"
	service_workspaces "getsturdy.com/api/pkg/workspaces/service"
	"getsturdy.com/api/vcs"
)

type LandRootResolver struct {
	landService      *service_land.Service
	workspaceService *service_workspaces.Service
	authService      *services_auth.Service

	workspaceResolver resolvers.WorkspaceRootResolver
}

func NewResolver(
	landService *service_land.Service,
	workspaceService *service_workspaces.Service,
	authService *services_auth.Service,

	workspaceResolver resolvers.WorkspaceRootResolver,
) resolvers.LandRootResovler {
	return &LandRootResolver{
		landService:       landService,
		workspaceService:  workspaceService,
		authService:       authService,
		workspaceResolver: workspaceResolver,
	}
}

func (r *LandRootResolver) LandWorkspaceChange(ctx context.Context, args resolvers.LandWorkspaceArgs) (resolvers.WorkspaceResolver, error) {
	ws, err := r.workspaceService.GetByID(ctx, string(args.Input.WorkspaceID))
	if err != nil {
		return nil, gqlerrors.Error(fmt.Errorf("failed to get workspace: %w", err))
	}

	if err := r.authService.CanWrite(ctx, ws); err != nil {
		return nil, gqlerrors.Error(err)
	}

	var diffOpts []vcs.DiffOption
	if args.Input.DiffMaxSize > 0 {
		diffOpts = append(diffOpts, vcs.WithGitMaxSize(args.Input.DiffMaxSize))
	}

	if _, err := r.landService.LandChange(ctx, ws, diffOpts...); err != nil {
		return nil, gqlerrors.Error(fmt.Errorf("failed to land change: %w", err))
	}

	return r.workspaceResolver.InternalWorkspace(ws), nil
}

func (r *LandRootResolver) PushWorkspace(ctx context.Context, args resolvers.PushWorkspaceArgs) (resolvers.WorkspaceResolver, error) {
	return nil, gqlerrors.ErrNotImplemented
}
