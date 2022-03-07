package open

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	eventsv2 "getsturdy.com/api/pkg/events/v2"
	"getsturdy.com/api/vcs"

	"getsturdy.com/api/pkg/snapshots"
	db3 "getsturdy.com/api/pkg/snapshots/db"
	"getsturdy.com/api/pkg/snapshots/snapshotter"
	vcs2 "getsturdy.com/api/pkg/snapshots/vcs"
	"getsturdy.com/api/pkg/view"
	"getsturdy.com/api/pkg/view/db"
	vcs_view "getsturdy.com/api/pkg/view/vcs"
	"getsturdy.com/api/pkg/workspaces"
	db_workspaces "getsturdy.com/api/pkg/workspaces/db"
	"getsturdy.com/api/vcs/executor"

	"go.uber.org/zap"
)

var ErrRebasing = errors.New("rebase in progress")

func OpenWorkspaceOnView(
	ctx context.Context,
	logger *zap.Logger,
	view *view.View,
	ws *workspaces.Workspace,
	viewRepo db.Repository,
	workspaceReader db_workspaces.WorkspaceReader,
	gitSnapshotter snapshotter.Snapshotter,
	snapshotRepo db3.Repository,
	workspaceWriter db_workspaces.WorkspaceWriter,
	executorProvider executor.Provider,
	eventSender *eventsv2.Publisher,
) error {
	if ws.ArchivedAt != nil {
		return fmt.Errorf("the workspace is archived")
	}
	if view.CodebaseID != ws.CodebaseID {
		return fmt.Errorf("the view and workspace does not belong to the same codebase")
	}
	if view.UserID != ws.UserID {
		return fmt.Errorf("the view and workspace does not belong to the same user")
	}

	// If the view that we're opening this workspace on has another workspace currently open, snapshot and save the changes
	currentWorkspaceOnView, err := workspaceReader.GetByViewID(view.ID, true)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return fmt.Errorf("could not find previous workspace on view")
	} else if err == nil {
		_, err := gitSnapshotter.Snapshot(currentWorkspaceOnView.CodebaseID, currentWorkspaceOnView.ID,
			snapshots.ActionPreCheckoutOtherWorkspace, snapshotter.WithOnView(*currentWorkspaceOnView.ViewID), snapshotter.WithMarkAsLatestInWorkspace())
		if errors.Is(err, snapshotter.ErrCantSnapshotRebasing) {
			return ErrRebasing
		}
		if err != nil {
			return fmt.Errorf("failed to snapshot: %w", err)
		}

		if err := workspaceWriter.UpdateFields(ctx, currentWorkspaceOnView.ID, db_workspaces.SetViewID(nil)); err != nil {
			return fmt.Errorf("failed to finalize previous workspace on view: %w", err)
		}
	}

	// If the workspace that we're opening on the view, currently is opened somewhere else, snapshot it and save the changes
	// The snapshot will also be used to "move" the contents to the new view
	if ws.ViewID != nil {
		_, err := gitSnapshotter.Snapshot(ws.CodebaseID, ws.ID, snapshots.ActionPreCheckoutOtherView,
			snapshotter.WithOnView(*ws.ViewID), snapshotter.WithMarkAsLatestInWorkspace())
		if err != nil {
			return fmt.Errorf("failed to snapshot: %w", err)
		}

		// TODO: unset view.workspace_id?
	}

	if err := executorProvider.New().
		Write(vcs_view.CheckoutBranch(ws.ID)).
		Write(func(repo vcs.RepoWriter) error {
			// Restore snapshot
			if ws.LatestSnapshotID != nil {
				snapshot, err := snapshotRepo.Get(*ws.LatestSnapshotID)
				if err != nil {
					return fmt.Errorf("failed to get snapshot: %w", err)
				}
				if err := vcs2.RestoreRepo(logger, repo, ws.CodebaseID, ws.ID, snapshot.ID, snapshot.CommitID); err != nil {
					return fmt.Errorf("failed to restore snapshot: %w", err)
				}
			} else {
				// vcs2.Restore does this as well, make sure to smudge files in this scenario as well
				if err := repo.LargeFilesPull(); err != nil {
					logger.Warn("failed to pull large files", zap.Error(err))
				}
			}
			return nil
		}).ExecView(ws.CodebaseID, view.ID, "OpenWorkspaceOnView"); err != nil {
		return fmt.Errorf("failed to open workspace on view: %w", err)
	}

	// Update workspace object
	ws.ViewID = &view.ID
	if err := workspaceWriter.UpdateFields(ctx, ws.ID, db_workspaces.SetViewID(&view.ID)); err != nil {
		return fmt.Errorf("failed to update workspace: %w", err)
	}

	// Update view object
	view.WorkspaceID = ws.ID
	if err := viewRepo.Update(view); err != nil {
		return fmt.Errorf("failed to update view obj: %w", err)
	}

	if err := eventSender.Workspace(ctx, ws.ID).ViewUpdated(view); err != nil {
		logger.Error("failed to send view updated event: %w", zap.Error(err))
		// do not fail
	}

	return nil
}
