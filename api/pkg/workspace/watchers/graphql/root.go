package graphql

import (
	"context"
	"fmt"

	"getsturdy.com/api/pkg/auth"
	service_auth "getsturdy.com/api/pkg/auth/service"
	gqlerrors "getsturdy.com/api/pkg/graphql/errors"
	"getsturdy.com/api/pkg/graphql/resolvers"
	"getsturdy.com/api/pkg/view/events"
	"getsturdy.com/api/pkg/workspace"
	service_workspace "getsturdy.com/api/pkg/workspace/service"
	service_workspace_watchers "getsturdy.com/api/pkg/workspace/watchers/service"

	"go.uber.org/zap"
)

type rootResolver struct {
	logger *zap.Logger

	workspaceWatcherService *service_workspace_watchers.Service
	workspaceService        service_workspace.Service

	authService *service_auth.Service

	eventsReader events.EventReader

	userRootResover       *resolvers.UserRootResolver
	workspaceRootResolver *resolvers.WorkspaceRootResolver
}

func NewRootResolver(
	logger *zap.Logger,

	workspaceWatcherService *service_workspace_watchers.Service,
	workspaceService service_workspace.Service,

	authService *service_auth.Service,

	eventsReader events.EventReader,

	userRootResover *resolvers.UserRootResolver,
	workspaceRootResolver *resolvers.WorkspaceRootResolver,
) resolvers.WorkspaceWatcherRootResolver {
	return &rootResolver{
		logger: logger.Named("workspaceWatchersRootResolver"),

		workspaceWatcherService: workspaceWatcherService,
		workspaceService:        workspaceService,

		authService: authService,

		eventsReader: eventsReader,

		userRootResover:       userRootResover,
		workspaceRootResolver: workspaceRootResolver,
	}
}

func (r *rootResolver) WatchWorkspace(ctx context.Context, args resolvers.WatchWorkspaceArgs) (resolvers.WorkspaceWatcherResolver, error) {
	userID, err := auth.UserID(ctx)
	if err != nil {
		return nil, gqlerrors.Error(err)
	}

	ws, err := r.workspaceService.GetByID(ctx, string(args.Input.WorkspaceID))
	if err != nil {
		return nil, gqlerrors.Error(fmt.Errorf("failed to fetch workspace: %w", err))
	}

	if err := r.authService.CanWrite(ctx, ws); err != nil {
		return nil, gqlerrors.Error(err)
	}

	watcher, err := r.workspaceWatcherService.Watch(ctx, userID, string(args.Input.WorkspaceID))
	if err != nil {
		return nil, gqlerrors.Error(err)
	}

	return &watcherResolver{
		Root:    r,
		Watcher: watcher,
	}, nil
}

func (r *rootResolver) UnwatchWorkspace(ctx context.Context, args resolvers.UnwatchWorkspaceArgs) (resolvers.WorkspaceWatcherResolver, error) {
	userID, err := auth.UserID(ctx)
	if err != nil {
		return nil, gqlerrors.Error(err)
	}

	ws, err := r.workspaceService.GetByID(ctx, string(args.Input.WorkspaceID))
	if err != nil {
		return nil, gqlerrors.Error(fmt.Errorf("failed to fetch workspace: %w", err))
	}

	if err := r.authService.CanWrite(ctx, ws); err != nil {
		return nil, gqlerrors.Error(err)
	}

	watcher, err := r.workspaceWatcherService.Unwatch(ctx, userID, string(args.Input.WorkspaceID))
	if err != nil {
		return nil, gqlerrors.Error(err)
	}

	return &watcherResolver{
		Root:    r,
		Watcher: watcher,
	}, nil
}

func (r *rootResolver) InternalWorkspaceWatchers(ctx context.Context, ws *workspace.Workspace) ([]resolvers.WorkspaceWatcherResolver, error) {
	watchers, err := r.workspaceWatcherService.ListWatchers(ctx, ws.ID)
	if err != nil {
		return nil, gqlerrors.Error(err)
	}

	resolvers := make([]resolvers.WorkspaceWatcherResolver, 0, len(watchers))
	for _, watcher := range watchers {
		resolvers = append(resolvers, &watcherResolver{
			Root:    r,
			Watcher: watcher,
		})
	}

	return resolvers, nil
}
