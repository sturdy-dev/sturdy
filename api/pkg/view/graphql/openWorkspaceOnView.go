package graphql

import (
	"context"
	"errors"

	gqlerrors "getsturdy.com/api/pkg/graphql/errors"
	"getsturdy.com/api/pkg/graphql/resolvers"
	"getsturdy.com/api/pkg/view/open"
)

func (r *ViewRootResolver) OpenWorkspaceOnView(ctx context.Context, args resolvers.OpenViewArgs) (resolvers.ViewResolver, error) {
	ws, err := r.workspaceReader.Get(string(args.Input.WorkspaceID))
	if err != nil {
		return nil, gqlerrors.Error(err, "workspace", "NotFound")
	}

	if err := r.authService.CanRead(ctx, ws); err != nil {
		return nil, gqlerrors.Error(err)
	}

	// Get view
	view, err := r.viewRepo.Get(string(args.Input.ViewID))
	if err != nil {
		return nil, gqlerrors.Error(err, "view", "NotFound")
	}

	if err := r.authService.CanWrite(ctx, view); err != nil {
		return nil, gqlerrors.Error(err)
	}

	// The workspace is already using this view
	if ws.ViewID != nil && view.ID == *ws.ViewID {
		// No-op
		return r.resolveView(ctx, args.Input.ViewID)
	}

	if view.UserID != ws.UserID {
		return r.CopyWorkspaceToView(ctx, resolvers.CopyViewArgs{Input: resolvers.CopyWorkspaceOnViewInput{
			ViewID:      args.Input.ViewID,
			WorkspaceID: args.Input.WorkspaceID,
		}})
	} else {
		if err := open.OpenWorkspaceOnView(r.logger, view, ws, r.viewRepo, r.workspaceReader,
			r.snapshotter, r.snapshotRepo, r.workspaceWriter, r.executorProvider, r.eventSender); errors.Is(err, open.ErrRebasing) {
			return nil, gqlerrors.Error(gqlerrors.ErrBadRequest, "message", "View is currently in rebasing state. Please resolve all the conflicts and try again.")
		} else if err != nil {
			return nil, gqlerrors.Error(err)
		}
	}

	return r.resolveView(ctx, args.Input.ViewID)
}
