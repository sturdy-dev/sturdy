package db

import (
	"context"
	"fmt"

	"getsturdy.com/api/pkg/licenses"
	"getsturdy.com/api/pkg/licenses/enterprise/cloud/validations"
	"github.com/jmoiron/sqlx"
)

type database struct {
	db *sqlx.DB
}

func New(db *sqlx.DB) Repository {
	return &database{db: db}
}

func (d *database) Create(ctx context.Context, validation *validations.Validation) error {
	if _, err := d.db.NamedExecContext(ctx, `
		INSERT INTO license_validations (
			id, license_id, timestamp, status
		) VALUES (
			:id, :license_id, :timestamp, :status
		)
	`, validation); err != nil {
		return fmt.Errorf("failed to create insert: %w", err)
	}
	return nil
}

func (d *database) ListLatest(ctx context.Context, licenseID licenses.ID) ([]*validations.Validation, error) {
	var validations []*validations.Validation
	if err := d.db.SelectContext(ctx, &validations, `
		SELECT
			id, license_id, timestamp, status
		FROM
			license_validations
		WHERE
			license_id = $1
		ORDER BY
			timestamp DESC
		LIMIT
			10
	`, licenseID); err != nil {
		return nil, fmt.Errorf("failed to select: %w", err)
	}
	return validations, nil
}
