package graphql

import (
	"context"
	"fmt"

	"getsturdy.com/api/pkg/auth"
	services_auth "getsturdy.com/api/pkg/auth/service"
	gqlerrors "getsturdy.com/api/pkg/graphql/errors"
	"getsturdy.com/api/pkg/graphql/resolvers"
	service_land "getsturdy.com/api/pkg/land/enterprise/service"
	service_users "getsturdy.com/api/pkg/users/service"
	service_workspaces "getsturdy.com/api/pkg/workspaces/service"
	"getsturdy.com/api/vcs"
)

type LandRootResolver struct {
	workspaceService *service_workspaces.Service
	landService      *service_land.Service
	authService      *services_auth.Service
	userService      service_users.Service

	workspaceResolver resolvers.WorkspaceRootResolver
}

func NewResolver(
	workspaceService *service_workspaces.Service,
	authService *services_auth.Service,
	landService *service_land.Service,
	userService service_users.Service,
	workspaceResolver resolvers.WorkspaceRootResolver,
) resolvers.LandRootResovler {
	return &LandRootResolver{
		workspaceService:  workspaceService,
		authService:       authService,
		userService:       userService,
		landService:       landService,
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
	ws, err := r.workspaceService.GetByID(ctx, string(args.Input.WorkspaceID))
	if err != nil {
		return nil, gqlerrors.Error(err)
	}

	if err := r.authService.CanWrite(ctx, ws); err != nil {
		return nil, gqlerrors.Error(err)
	}

	userID, err := auth.UserID(ctx)
	if err != nil {
		return nil, gqlerrors.Error(err)
	}

	if args.Input.LandOnSturdyAndPushTracked != nil && *args.Input.LandOnSturdyAndPushTracked {
		if err := r.landService.LandOnSturdyAndPushTracked(ctx, ws); err != nil {
			return nil, gqlerrors.Error(err)
		}
	} else {
		user, err := r.userService.GetByID(ctx, userID)
		if err != nil {
			return nil, gqlerrors.Error(err)
		}

		if err := r.landService.Push(ctx, user, ws); err != nil {
			return nil, gqlerrors.Error(err)
		}
	}

	return r.workspaceResolver.InternalWorkspace(ws), nil
}
