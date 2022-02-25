package watchers

import (
	"time"

	"getsturdy.com/api/pkg/users"
)

type Status string

const (
	StatusUndefined Status = "undefined"
	StatusWatching  Status = "watching"
	StatusIgnored   Status = "ignored"
)

type Watcher struct {
	WorkspaceID string    `db:"workspace_id"`
	UserID      users.ID  `db:"user_id"`
	Status      Status    `db:"status"`
	CreatedAt   time.Time `db:"created_at"`
}
