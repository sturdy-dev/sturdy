package graphql

import (
	"context"
	"database/sql"
	"errors"

	gqlerrors "getsturdy.com/api/pkg/graphql/errors"
	"getsturdy.com/api/pkg/graphql/resolvers"
)

func (r *WorkspaceRootResolver) RemovePatches(ctx context.Context, args resolvers.RemovePatchesArgs) (resolvers.WorkspaceResolver, error) {
	ws, err := r.workspaceReader.Get(string(args.Input.WorkspaceID))
	if err != nil {
		return nil, gqlerrors.Error(err)
	}

	if err := r.authService.CanWrite(ctx, ws); err != nil {
		return nil, gqlerrors.Error(err)
	}

	suggestion, err := r.suggestionsService.GetByWorkspaceID(ctx, ws.ID)
	switch {
	case err == nil:
		if err := r.suggestionsService.RemovePatches(ctx, suggestion, args.Input.HunkIDs...); err != nil {
			return nil, gqlerrors.Error(err)
		}
	case errors.Is(err, sql.ErrNoRows):
		if err := r.workspaceService.RemovePatches(ctx, ws, args.Input.HunkIDs...); err != nil {
			return nil, gqlerrors.Error(err)
		}
	default:
		return nil, gqlerrors.Error(err)
	}

	return &WorkspaceResolver{root: r, w: ws}, nil
}
