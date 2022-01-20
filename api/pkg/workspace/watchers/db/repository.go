package db

import (
	"context"

	"getsturdy.com/api/pkg/workspace/watchers"
)

type Repository interface {
	// Create creates a new watcher in the database.
	Create(context.Context, *watchers.Watcher) error
	// ListWatchingByWorkspaceID returns a list of watchers watching the given workspace.
	ListWatchingByWorkspaceID(context.Context, string) ([]*watchers.Watcher, error)
	// GetByUserIDWorkspaceID returns a watcher by user ID and workspace ID.
	GetByUserIDAndWorkspaceID(context.Context, string, string) (*watchers.Watcher, error)
}
