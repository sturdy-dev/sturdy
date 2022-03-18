package gc

import (
	"time"

	"getsturdy.com/api/pkg/codebases"
)

type CodebaseGarbageStatus struct {
	CodebaseID     codebases.ID `db:"codebase_id"`
	CompletedAt    time.Time    `db:"completed_at"`
	DurationMillis int64        `db:"duration_millis"`
}
