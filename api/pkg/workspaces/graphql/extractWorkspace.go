package graphql

import (
	"context"
	"fmt"

	"getsturdy.com/api/pkg/auth"
	"getsturdy.com/api/pkg/events"
	gqlerrors "getsturdy.com/api/pkg/graphql/errors"
	"getsturdy.com/api/pkg/graphql/resolvers"
	"getsturdy.com/api/pkg/workspaces"
	service_workspace "getsturdy.com/api/pkg/workspaces/service"

	"github.com/graph-gophers/graphql-go"
	"go.uber.org/zap"
)

func (r *WorkspaceRootResolver) ExtractWorkspace(ctx context.Context, args resolvers.ExtractWorkspaceArgs) (resolvers.WorkspaceResolver, error) {
	ws, err := r.workspaceReader.Get(string(args.Input.WorkspaceID))
	if err != nil {
		return nil, gqlerrors.Error(err)
	}

	if err := r.authService.CanWrite(ctx, ws); err != nil {
		return nil, gqlerrors.Error(err)
	}

	extractor := &workspaceExtractor{
		workspaceService: r.workspaceService,
	}

	newWorkspace, err := extractor.Extract(ctx, ws, args.Input.PatchIDs)
	if err != nil {
		r.logger.Error("filed to extract workspace", zap.Error(err))
		return nil, gqlerrors.Error(err)
	}

	if err := r.eventsSender.Codebase(newWorkspace.CodebaseID, events.CodebaseUpdated, newWorkspace.CodebaseID.String()); err != nil {
		r.logger.Error("failed to send codebase event", zap.Error(err))
	}

	return r.Workspace(ctx, resolvers.WorkspaceArgs{ID: graphql.ID(newWorkspace.ID)})
}

type workspaceExtractor struct {
	workspaceService *service_workspace.Service
}

func (r *workspaceExtractor) Extract(ctx context.Context, src *workspaces.Workspace, patchIDs []string) (*workspaces.Workspace, error) {
	userID, err := auth.UserID(ctx)
	if err != nil {
		return nil, fmt.Errorf("could not get authed user: %w", err)
	}

	name := ""
	if src.Name != nil {
		name = fmt.Sprintf("Fork of %s", *src.Name)
	}

	dist, err := r.workspaceService.CreateFromWorkspace(ctx, src, userID, name)
	if err != nil {
		return nil, fmt.Errorf("failed to copy a workspace: %w", err)
	}

	if err := r.workspaceService.CopyPatches(ctx, dist, src, service_workspace.WithPatchIDs(patchIDs)); err != nil {
		return nil, fmt.Errorf("failed to copy patches: %w", err)
	}

	return dist, nil
}
