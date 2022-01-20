package watchers

import "time"

type Status string

const (
	StatusUndefined Status = "undefined"
	StatusWatching  Status = "watching"
	StatusIgnored   Status = "ignored"
)

type Watcher struct {
	WorkspaceID string    `db:"workspace_id"`
	UserID      string    `db:"user_id"`
	Status      Status    `db:"status"`
	CreatedAt   time.Time `db:"created_at"`
}
