package db

import (
	"context"
	"fmt"
	"getsturdy.com/api/pkg/integrations/buildkite"

	"github.com/jmoiron/sqlx"
)

var _ Repository = &database{}

type database struct {
	db *sqlx.DB
}

func NewDatabase(db *sqlx.DB) Repository {
	return &database{db: db}
}

func (d *database) Create(ctx context.Context, cfg *buildkite.Config) error {
	if _, err := d.db.NamedExecContext(ctx, `
		INSERT INTO ci_configurations_buildkite 
			(id, codebase_id, integration_id, organization_name, pipeline_name, api_token, webhook_secret, created_at)
		VALUES
			(:id, :codebase_id, :integration_id, :organization_name, :pipeline_name, :api_token, :webhook_secret, :created_at)
	`, cfg); err != nil {
		return fmt.Errorf("failed to insert: %w", err)
	}
	return nil
}

func (d *database) Update(ctx context.Context, cfg *buildkite.Config) error {
	if _, err := d.db.NamedExecContext(ctx, `
		UPDATE ci_configurations_buildkite 
		SET
			organization_name = :organization_name,
			pipeline_name = :pipeline_name,
			api_token = :api_token,
			webhook_secret = :webhook_secret, 
			updated_at = :updated_at
		WHERE
			id = :id
	`, cfg); err != nil {
		return fmt.Errorf("failed to update: %w", err)
	}
	return nil
}

func (d *database) GetConfigsByCodebaseID(ctx context.Context, codebaseID string) ([]*buildkite.Config, error) {
	var cfgs []*buildkite.Config
	if err := d.db.SelectContext(ctx, &cfgs, `
		SELECT 
			id, codebase_id, integration_id, organization_name, pipeline_name, api_token, webhook_secret, created_at
		FROM ci_configurations_buildkite
		WHERE codebase_id = $1
	`, codebaseID); err != nil {
		return nil, fmt.Errorf("failed to get config: %w", err)
	}
	return cfgs, nil
}

func (d *database) GetConfigByIntegrationID(ctx context.Context, integrationID string) (*buildkite.Config, error) {
	var cfg buildkite.Config
	if err := d.db.GetContext(ctx, &cfg, `
		SELECT 
			id, codebase_id, integration_id, organization_name, pipeline_name, api_token, webhook_secret, created_at
		FROM ci_configurations_buildkite
		WHERE integration_id = $1
	`, integrationID); err != nil {
		return nil, fmt.Errorf("failed to get config: %w", err)
	}
	return &cfg, nil
}
