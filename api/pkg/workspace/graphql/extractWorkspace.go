package graphql

import (
	"context"
	"fmt"

	"getsturdy.com/api/pkg/events"
	gqlerrors "getsturdy.com/api/pkg/graphql/errors"
	"getsturdy.com/api/pkg/graphql/resolvers"
	"getsturdy.com/api/pkg/snapshots/snapshotter"
	"getsturdy.com/api/pkg/workspace"
	db_workspace "getsturdy.com/api/pkg/workspace/db"
	service_workspace "getsturdy.com/api/pkg/workspace/service"

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
		workspaceWriter:  r.workspaceWriter,
		gitSnapshotter:   r.gitSnapshotter,
		workspaceService: r.workspaceService,
	}

	newWorkspace, err := extractor.Extract(ctx, ws, args.Input.PatchIDs)
	if err != nil {
		r.logger.Error("filed to extract workspace", zap.Error(err))
		return nil, gqlerrors.Error(err)
	}

	if err := r.eventsSender.Codebase(newWorkspace.CodebaseID, events.CodebaseUpdated, newWorkspace.CodebaseID); err != nil {
		r.logger.Error("failed to send codebase event", zap.Error(err))
	}

	return r.Workspace(ctx, resolvers.WorkspaceArgs{ID: graphql.ID(newWorkspace.ID)})
}

type workspaceExtractor struct {
	workspaceWriter  db_workspace.WorkspaceWriter
	gitSnapshotter   snapshotter.Snapshotter
	workspaceService service_workspace.Service
}

func (r *workspaceExtractor) Extract(ctx context.Context, src *workspace.Workspace, patchIDs []string) (*workspace.Workspace, error) {
	dist, err := r.copyWorkspace(ctx, src)
	if err != nil {
		return nil, fmt.Errorf("failed to copy a workspace: %w", err)
	}

	if err := r.workspaceService.CopyPatches(ctx, dist, src, service_workspace.WithPatchIDs(patchIDs)); err != nil {
		return nil, fmt.Errorf("failed to copy patches: %w", err)
	}

	return dist, nil
}

func (r *workspaceExtractor) copyWorkspace(ctx context.Context, from *workspace.Workspace) (*workspace.Workspace, error) {
	changeID := ""
	if from.HeadCommitID != nil {
		changeID = *from.HeadCommitID
	}

	name := ""
	if from.Name != nil {
		name = fmt.Sprintf("Fork of %s", *from.Name)
	}

	createRequest := service_workspace.CreateWorkspaceRequest{
		UserID:     from.UserID,
		CodebaseID: from.CodebaseID,
		Name:       name,
		ChangeID:   changeID,
	}

	newWorkspace, err := r.workspaceService.Create(ctx, createRequest)
	if err != nil {
		return nil, fmt.Errorf("faliled to create a workspace")
	}
	return newWorkspace, nil
}
