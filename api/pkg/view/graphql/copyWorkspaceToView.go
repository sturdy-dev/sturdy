package graphql

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	gqlerrors "getsturdy.com/api/pkg/graphql/errors"
	"getsturdy.com/api/pkg/graphql/resolvers"
	"getsturdy.com/api/pkg/snapshots"
	"getsturdy.com/api/pkg/snapshots/snapshotter"
	"getsturdy.com/api/pkg/view/vcs"
	vcsvcs "getsturdy.com/api/vcs"

	"go.uber.org/zap"
)

func (r *ViewRootResolver) CopyWorkspaceToView(ctx context.Context, args resolvers.CopyViewArgs) (resolvers.ViewResolver, error) {
	view, err := r.viewRepo.Get(string(args.Input.ViewID))
	if err != nil {
		return nil, gqlerrors.Error(err)
	}

	if err := r.authService.CanWrite(ctx, view); err != nil {
		return nil, gqlerrors.Error(err)
	}

	ws, err := r.workspaceReader.Get(string(args.Input.WorkspaceID))
	if err != nil {
		return nil, gqlerrors.Error(err, "workspace", "NotFound")
	}

	if err := r.authService.CanRead(ctx, ws); err != nil {
		return nil, gqlerrors.Error(err)
	}

	if ws.ArchivedAt != nil {
		return nil, gqlerrors.Error(gqlerrors.ErrBadRequest, "message", "the workspace is archived")
	}
	if view.CodebaseID != ws.CodebaseID {
		return nil, gqlerrors.Error(gqlerrors.ErrBadRequest, "message", "the view and workspace does not belong to the same codebase")
	}

	// If the view that we're opening this workspace on has another workspace currently open, snapshot and save the changes
	currentWorkspaceOnView, err := r.workspaceReader.GetByViewID(view.ID, true)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return nil, fmt.Errorf("could not find previous workspace on view")
	} else if err == nil {
		snapshot, err := r.snapshotter.Snapshot(currentWorkspaceOnView.CodebaseID, currentWorkspaceOnView.ID, snapshots.ActionPreCheckoutOtherWorkspace, snapshotter.WithOnView(*currentWorkspaceOnView.ViewID))
		if err != nil {
			return nil, gqlerrors.Error(fmt.Errorf("failed to snapshot: %w", err))
		}
		currentWorkspaceOnView.LatestSnapshotID = &snapshot.ID
		currentWorkspaceOnView.ViewID = nil // The workspace no longer has any view open
		if err := r.workspaceWriter.Update(currentWorkspaceOnView); err != nil {
			return nil, gqlerrors.Error(fmt.Errorf("failed to finalize previous workspace on view: %w", err))
		}
	}

	// TODO: Create a snapshot of the authorative view (or use latest snapshot)
	if err := r.executorProvider.New().Write(func(repo vcsvcs.RepoWriter) error {
		// For backwards compatability? Skip checkout if already checked out
		headBranch, err := repo.HeadBranch()
		if err != nil {
			return fmt.Errorf("failed to open repo get branch")
		}
		if headBranch != ws.ID {
			// Checkout
			if err := vcs.SetWorkspaceRepo(repo, ws.ID); err != nil {
				return fmt.Errorf("failed to checkout view: %w", err)
			}
		}
		return nil
	}).ExecView(ws.CodebaseID, view.ID, "copyWorkspaceToView"); err != nil {
		return nil, gqlerrors.Error(err)
	}

	// Update view object
	// Save that this view is operating on a copy
	view.WorkspaceID = ws.ID
	if err := r.viewRepo.Update(view); err != nil {
		return nil, fmt.Errorf("failed to update view obj: %w", err)
	}

	if _, err := r.workspaceWatchersService.Watch(ctx, view.UserID, ws.ID); err != nil {
		r.logger.Error("failed to watch workspace", zap.Error(err))
		// do not fail
	}

	return r.resolveView(ctx, args.Input.ViewID)
}
