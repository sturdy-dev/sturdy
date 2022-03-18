package integrations

import (
	"time"

	"getsturdy.com/api/pkg/codebases"
	"getsturdy.com/api/pkg/integrations/providers"
)

type Integration struct {
	ID           string                 `db:"id"`
	CodebaseID   codebases.ID           `db:"codebase_id"`
	Provider     providers.ProviderName `db:"provider"`
	ProviderType providers.ProviderType `db:"provider_type"`
	SeedFiles    []string               `db:"seed_files"` // TODO: Can we move or remove this? It's specific to Build providers
	CreatedAt    time.Time              `db:"created_at"`
	UpdatedAt    time.Time              `db:"updated_at"`
	DeletedAt    *time.Time             `db:"deleted_at"`
}
