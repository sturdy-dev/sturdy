package graphql

import (
	"context"
	"fmt"

	service_auth "getsturdy.com/api/pkg/auth/service"
	"getsturdy.com/api/pkg/codebases"
	service_codebase "getsturdy.com/api/pkg/codebases/service"
	gqlerror "getsturdy.com/api/pkg/graphql/errors"
	"getsturdy.com/api/pkg/graphql/resolvers"
	"getsturdy.com/api/pkg/remote/enterprise/service"
	service_user "getsturdy.com/api/pkg/users/service"
	service_workspace "getsturdy.com/api/pkg/workspaces/service"
)

type remoteRootResolver struct {
	service          *service.Service
	workspaceService service_workspace.Service
	authService      *service_auth.Service
	codebaseService  *service_codebase.Service
	userService      service_user.Service
}

func New(
	service *service.Service,
	workspaceService service_workspace.Service,
	authService *service_auth.Service,
	codebaseService *service_codebase.Service,
	userService service_user.Service,
) resolvers.RemoteRootResolver {
	return &remoteRootResolver{
		service:          service,
		workspaceService: workspaceService,
		authService:      authService,
		codebaseService:  codebaseService,
		userService:      userService,
	}
}

func (r *remoteRootResolver) InternalRemoteByCodebaseID(ctx context.Context, codebaseID codebases.ID) (resolvers.RemoteResolver, error) {
	cb, err := r.codebaseService.GetByID(ctx, codebaseID)
	if err != nil {
		return nil, gqlerror.Error(err)
	}

	if err := r.authService.CanWrite(ctx, cb); err != nil {
		return nil, gqlerror.Error(err)
	}

	rem, err := r.service.Get(ctx, codebaseID)
	if err != nil {
		return nil, gqlerror.Error(err)
	}

	return &resolver{remote: rem}, nil
}

func (r *remoteRootResolver) CreateOrUpdateCodebaseRemote(ctx context.Context, args resolvers.CreateOrUpdateCodebaseRemoteArgsArgs) (resolvers.RemoteResolver, error) {
	codebaseID := codebases.ID(args.Input.CodebaseID)
	cb, err := r.codebaseService.GetByID(ctx, codebaseID)
	if err != nil {
		return nil, gqlerror.Error(err)
	}

	if err := r.authService.CanWrite(ctx, cb); err != nil {
		return nil, gqlerror.Error(err)
	}

	rem, err := r.service.SetRemote(
		ctx,
		codebaseID,
		args.Input.Name,
		args.Input.Url,
		args.Input.BasicAuthUsername,
		args.Input.BasicAuthPassword,
		args.Input.TrackedBranch,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to add remote: %w", err)
	}

	return &resolver{remote: rem}, nil
}

func (r *remoteRootResolver) PushWorkspaceToRemote(ctx context.Context, args resolvers.PushWorkspaceToRemoteArgs) (resolvers.WorkspaceResolver, error) {
	return nil, gqlerror.ErrNotImplemented
	/*
		ws, err := r.workspaceService.GetByID(ctx, string(args.Input.WorkspaceID))
		if err != nil {
			return nil, gqlerror.Error(fmt.Errorf("could not get workspace: %w", err))
		}

		if err := r.authService.CanWrite(ctx, ws); err != nil {
			return nil, gqlerror.Error(err)
		}

		userID, err := auth.UserID(ctx)
		if err != nil {
			return nil, gqlerror.Error(err)
		}

		user, err := r.userService.GetByID(ctx, userID)
		if err != nil {
			return nil, gqlerror.Error(err)
		}

		if err := r.service.Push(ctx, user, ws); err != nil {
			return nil, gqlerror.Error(fmt.Errorf("failed to push workspace: %w", err))
		}

		return r.workspaceRootResolver.Workspace(ctx, resolvers.WorkspaceArgs{ID: args.Input.WorkspaceID})
	*/
}

func (r *remoteRootResolver) RemoteFetchToTrunk(ctx context.Context, args resolvers.PushWorkspaceToRemoteArgs) (resolvers.WorkspaceResolver, error) {
	return nil, gqlerror.ErrNotImplemented
	/*
		ws, err := r.workspaceService.GetByID(ctx, string(args.Input.WorkspaceID))
		if err != nil {
			return nil, gqlerror.Error(fmt.Errorf("could not get workspace: %w", err))
		}

		if err := r.authService.CanWrite(ctx, ws); err != nil {
			return nil, gqlerror.Error(err)
		}

		if err := r.service.Pull(ctx, ws.CodebaseID); err != nil {
			return nil, gqlerror.Error(fmt.Errorf("failed to pull remote to trunk: %w", err))
		}

		return r.workspaceRootResolver.Workspace(ctx, resolvers.WorkspaceArgs{ID: args.Input.WorkspaceID})
	*/
}

func (r *remoteRootResolver) FetchRemoteToTrunk(ctx context.Context, args resolvers.FetchRemoteToTrunkArgs) (resolvers.CodebaseResolver, error) {
	return nil, gqlerror.ErrNotImplemented
	/*
		cb, err := r.codebaseService.GetByID(ctx, string(args.Input.CodebaseID))
		if err != nil {
			return nil, gqlerror.Error(fmt.Errorf("could not get codebase: %w", err))
		}

		if err := r.authService.CanWrite(ctx, cb); err != nil {
			return nil, gqlerror.Error(err)
		}

		if err := r.service.Pull(ctx, cb.ID); err != nil {
			return nil, gqlerror.Error(fmt.Errorf("failed to fetch codebase: %w", err))
		}

		return r.codebaseRootResolver.Codebase(ctx, resolvers.CodebaseArgs{ID: &args.Input.CodebaseID})
	*/
}
