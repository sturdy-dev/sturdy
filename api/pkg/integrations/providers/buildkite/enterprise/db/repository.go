package db

import (
	"context"

	"getsturdy.com/api/pkg/codebases"
	"getsturdy.com/api/pkg/integrations/providers/buildkite"
)

type Repository interface {
	Create(context.Context, *buildkite.Config) error
	Update(context.Context, *buildkite.Config) error
	GetConfigsByCodebaseID(context.Context, codebases.ID) ([]*buildkite.Config, error)
	GetConfigByIntegrationID(ctx context.Context, integrationID string) (*buildkite.Config, error)
}
