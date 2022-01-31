package db

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"getsturdy.com/api/pkg/licenses"
	"github.com/jmoiron/sqlx"
)

type database struct {
	db *sqlx.DB
}

func New(db *sqlx.DB) Repository {
	return &database{
		db: db,
	}
}

func (d *database) Create(ctx context.Context, license *licenses.License) error {
	if _, err := d.db.NamedExecContext(ctx, `
		INSERT INTO licenses (
			id,
			organization_id,
			key,
			created_at,
			expires_at,
			seats
		) VALUES (
			:id,
			:organization_id,
			:key,
			:created_at,
			:expires_at,
			:seats
		)
	`, license); err != nil {
		return fmt.Errorf("failed to insert: %w", err)
	}
	return nil
}

func (d *database) Get(ctx context.Context, id licenses.ID) (*licenses.License, error) {
	license := &licenses.License{}
	if err := d.db.GetContext(ctx, license, `
		SELECT
			id,
			organization_id,
			key,
			created_at,
			expires_at,
			seats
		FROM licenses
		WHERE id = $1
	`, id); errors.Is(err, sql.ErrNoRows) {
		return nil, ErrNotFound
	} else if err != nil {
		return nil, fmt.Errorf("failed to select: %w", err)
	}
	return license, nil
}

func (d *database) GetByKey(ctx context.Context, key string) (*licenses.License, error) {
	license := &licenses.License{}
	if err := d.db.GetContext(ctx, license, `
		SELECT
			id,
			organization_id,
			key,
			created_at,
			expires_at,
			seats
		FROM licenses
		WHERE key = $1
	`, key); errors.Is(err, sql.ErrNoRows) {
		return nil, ErrNotFound
	} else if err != nil {
		return nil, fmt.Errorf("failed to select: %w", err)
	}
	return license, nil
}

func (d *database) ListByOrganizationID(ctx context.Context, oranizationID string) ([]*licenses.License, error) {
	licenses := []*licenses.License{}
	if err := d.db.SelectContext(ctx, &licenses, `
		SELECT
			id,
			organization_id,
			key,
			created_at,
			expires_at,
			seats
		FROM licenses
		WHERE organization_id = $1
	`, oranizationID); err != nil {
		return nil, fmt.Errorf("failed to select: %w", err)
	}
	return licenses, nil
}
