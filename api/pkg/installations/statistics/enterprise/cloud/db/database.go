package db

import (
	"context"
	"database/sql"
	"errors"
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

func (d *database) GetByLicenseKey(ctx context.Context, key string) (*statistics.Statistic, error) {
	statistic := &statistics.Statistic{}
	if err := d.db.GetContext(ctx, statistic, `
		SELECT
			installation_id,
			license_key,
			version,
			ip,
			recorded_at,
			received_at,
			users_count,
			codebases_count
		FROM 
			installation_statistics
		WHERE 
			license_key = $1
		ORDER BY 
			recorded_at DESC
		LIMIT 1
	`, key); errors.Is(err, sql.ErrNoRows) {
		return nil, ErrNotFound
	} else if err != nil {
		return nil, fmt.Errorf("failed to get: %w", err)
	}
	return statistic, nil
}
