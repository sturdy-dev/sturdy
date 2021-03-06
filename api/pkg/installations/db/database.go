package db

import (
	"context"
	"fmt"

	"getsturdy.com/api/pkg/installations"

	"github.com/jmoiron/sqlx"
)

type Database struct {
	db *sqlx.DB
}

func NewDatabase(db *sqlx.DB) Repository {
	return &Database{db: db}
}

func (d *Database) Create(ctx context.Context, installation *installations.Installation) error {
	if _, err := d.db.NamedExecContext(ctx, `
		INSERT INTO installations (id, license_key) VALUES (:id, :license_key)
	`, installation); err != nil {
		return fmt.Errorf("failed to insert: %w", err)
	}
	return nil
}

func (d *Database) ListAll(ctx context.Context) ([]*installations.Installation, error) {
	var list []*installations.Installation
	if err := d.db.SelectContext(ctx, &list, `
		SELECT id, license_key FROM installations
	`); err != nil {
		return nil, fmt.Errorf("failed to select: %w", err)
	}
	return list, nil
}

func (d *Database) Update(ctx context.Context, installation *installations.Installation) error {
	_, err := d.db.NamedExecContext(ctx, `UPDATE installations
		SET license_key = :license_key
		WHERE id = :id`, installation)
	if err != nil {
		return fmt.Errorf("failed to update installation: %w", err)
	}
	return nil
}
