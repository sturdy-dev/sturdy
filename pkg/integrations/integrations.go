package integrations

import (
	"time"
)

type Integration struct {
	ID         string       `db:"id"`
	CodebaseID string       `db:"codebase_id"`
	Provider   ProviderType `db:"provider"`
	SeedFiles  []string     `db:"seed_files"`
	CreatedAt  time.Time    `db:"created_at"`
	UpdatedAt  time.Time    `db:"updated_at"`
	DeletedAt  *time.Time   `db:"deleted_at"`
}
