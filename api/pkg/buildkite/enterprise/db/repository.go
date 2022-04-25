package db

import (
	"context"

	"getsturdy.com/api/pkg/buildkite"
	"getsturdy.com/api/pkg/codebases"
)

type Repository interface {
	Create(context.Context, *buildkite.Config) error
	Update(context.Context, *buildkite.Config) error
	GetConfigsByCodebaseID(context.Context, codebases.ID) ([]*buildkite.Config, error)
	GetConfigByIntegrationID(ctx context.Context, integrationID string) (*buildkite.Config, error)
}
