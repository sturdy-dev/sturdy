package db

import (
	"context"
	"fmt"

	"getsturdy.com/api/pkg/installations/statistics"
	"github.com/jmoiron/sqlx"
)

type database struct {
	db *sqlx.DB
}

func New(db *sqlx.DB) Repository {
	return &database{db}
}

func (d *database) Create(ctx context.Context, statistics *statistics.Statistic) error {
	if _, err := d.db.NamedExecContext(ctx, `
		INSERT INTO installation_statistics (
			installation_id,
			license_key,
			version,
			ip,
			recorded_at,
			received_at,
			users_count,
			codebases_count
		) VALUES (
			:installation_id,
			:license_key,
			:version,
			:ip,
			:recorded_at,
			:received_at,
			:users_count,
			:codebases_count
		)
	`, statistics); err != nil {
		return fmt.Errorf("failed to insert: %w", err)
	}
	return nil
}
