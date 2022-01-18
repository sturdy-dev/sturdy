package gc

import "time"

type CodebaseGarbageStatus struct {
	CodebaseID     string    `db:"codebase_id"`
	CompletedAt    time.Time `db:"completed_at"`
	DurationMillis int64     `db:"duration_millis"`
}
