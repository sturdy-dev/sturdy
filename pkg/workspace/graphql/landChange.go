package graphql

import (
	"context"
	"fmt"

	gqlerrors "mash/pkg/graphql/errors"
	"mash/pkg/graphql/resolvers"
	"mash/vcs"
)

func (r *WorkspaceRootResolver) LandWorkspaceChange(ctx context.Context, args resolvers.LandWorkspaceArgs) (resolvers.WorkspaceResolver, error) {
	ws, err := r.workspaceReader.Get(string(args.Input.WorkspaceID))
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

	if err := r.workspaceService.LandChange(ctx, ws, args.Input.PatchIDs, diffOpts...); err != nil {
		return nil, gqlerrors.Error(fmt.Errorf("failed to land change: %w", err))
	}

	return r.Workspace(ctx, resolvers.WorkspaceArgs{ID: args.Input.WorkspaceID})
}
