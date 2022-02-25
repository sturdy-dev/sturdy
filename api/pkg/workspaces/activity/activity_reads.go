package activity

import (
	"time"

	"getsturdy.com/api/pkg/users"
)

type WorkspaceActivityReads struct {
	ID                string    `db:"id"`
	UserID            users.ID  `db:"user_id"`
	WorkspaceID       string    `db:"workspace_id"`
	LastReadCreatedAt time.Time `db:"last_read_created_at"`
}
