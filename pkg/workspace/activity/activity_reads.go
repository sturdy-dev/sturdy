package activity

import "time"

type WorkspaceActivityReads struct {
	ID                string    `db:"id"`
	UserID            string    `db:"user_id"`
	WorkspaceID       string    `db:"workspace_id"`
	LastReadCreatedAt time.Time `db:"last_read_created_at"`
}
