package view

import (
	"time"

	"getsturdy.com/api/pkg/author"
	"getsturdy.com/api/pkg/jsontime"
	"getsturdy.com/api/pkg/users"
)

type View struct {
	ID         string   `db:"id" json:"id"`
	UserID     users.ID `db:"user_id" json:"user_id"`
	CodebaseID string   `db:"codebase_id" json:"codebase_id"`

	// Deprecated: in favour for workspace.ViewID
	// TODO: Make nulllable, migrate, and delete?
	WorkspaceID string `db:"workspace_id" json:"workspace_id"`

	// Deprecated: use MountPath and MountHostname instead
	Name *string `db:"name" json:"name"`

	MountPath     *string `db:"mount_path" json:"mount_path"`
	MountHostname *string `db:"mount_hostname" json:"mount_hostname"`

	// When the view was last used by a fuse-client
	LastUsedAt *time.Time `db:"last_used_at" json:"last_used_at"`

	CreatedAt *time.Time `db:"created_at" json:"created_at"`
}

type ViewJSON struct {
	View
	LastUsedAt jsontime.Time `json:"last_used_at"`
}

type ViewWithMetadataJSON struct {
	ViewJSON
	ViewWorkspaceMeta ViewWorkspaceMeta `json:"workspace"`
	User              author.Author     `json:"user"`
}

type ViewWorkspaceMeta struct {
	ID   string  `json:"id"`
	Name *string `json:"name"`
}

type ViewWorkspaceSnapshot struct {
	ID          string     `db:"id"`
	ViewID      string     `db:"view_id"`
	WorkspaceID string     `db:"workspace_id"`
	SnapshotID  string     `db:"snapshot_id"`
	CreatedAt   time.Time  `db:"created_at"`
	UpdatedAt   *time.Time `db:"updated_at"`
}
