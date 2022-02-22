package meta

import (
	"context"
	"fmt"

	"getsturdy.com/api/pkg/events"
	"getsturdy.com/api/pkg/snapshots"
	worker_snapshotter "getsturdy.com/api/pkg/snapshots/worker"
	"getsturdy.com/api/pkg/view"
	db_workspaces "getsturdy.com/api/pkg/workspaces/db"
	workspace_meta "getsturdy.com/api/pkg/workspaces/meta"
)

type ViewUpdatedFunc func(ctx context.Context, view *view.View, action snapshots.Action) error

// NewViewUpdatedFunc returns a function that sends events for updates of a views
func NewViewUpdatedFunc(
	workspaceReader db_workspaces.WorkspaceReader,
	workspaceWriter db_workspaces.WorkspaceWriter,
	eventsSender events.EventSender,
	snapshotterQueue worker_snapshotter.Queue,
) ViewUpdatedFunc {
	return func(ctx context.Context, view *view.View, action snapshots.Action) error {
		// Workspace has updated
		if err := workspace_meta.Updated(ctx, workspaceReader, workspaceWriter, view.WorkspaceID); err != nil {
			return fmt.Errorf("error updating workspace meta: %w", err)
		}

		// Add to snapshotter queue
		if err := snapshotterQueue.Enqueue(ctx, view.CodebaseID, view.ID, view.WorkspaceID, []string{"."}, action); err != nil {
			return fmt.Errorf("failed to enqueue snapshot: %w", err)
		}

		if err := eventsSender.Codebase(view.CodebaseID, events.ViewUpdated, view.ID); err != nil {
			return fmt.Errorf("failed to send view updated event: %w", err)
		}

		return nil
	}
}
