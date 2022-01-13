package db

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"

	"mash/pkg/organization"
)

type Repository interface {
	Get(ctx context.Context, id string) (*organization.Organization, error)
	GetFirst(ctx context.Context) (*organization.Organization, error)
	Create(ctx context.Context, org organization.Organization) error
	Update(ctx context.Context, org *organization.Organization) error
}

type repository struct {
	db *sqlx.DB
}

func New(db *sqlx.DB) Repository {
	return &repository{db: db}
}

func (r *repository) GetFirst(ctx context.Context) (*organization.Organization, error) {
	var org organization.Organization
	if err := r.db.GetContext(ctx, &org, `SELECT id, name, created_at, deleted_at FROM organizations`); err != nil {
		return nil, fmt.Errorf("could not get organization: %w", err)
	}
	return &org, nil
}

func (r *repository) Get(ctx context.Context, id string) (*organization.Organization, error) {
	var org organization.Organization
	if err := r.db.GetContext(ctx, &org, `SELECT id, name, created_at, deleted_at FROM organizations WHERE id = $1`, id); err != nil {
		return nil, fmt.Errorf("could not get organization: %w", err)
	}
	return &org, nil
}

func (r *repository) Create(ctx context.Context, org organization.Organization) error {
	if _, err := r.db.NamedExecContext(ctx, `INSERT INTO organizations (id, name, created_at, created_by, deleted_at, deleted_by) VALUES (:id, :name, :created_at, :created_by, :deleted_at, :deleted_by)`, org); err != nil {
		return fmt.Errorf("failed to create organization: %w", err)
	}
	return nil
}

func (r *repository) Update(ctx context.Context, org *organization.Organization) error {
	if _, err := r.db.NamedExecContext(ctx, `UPDATE organizations
		SET name = :name,
	    	deleted_at = :deleted_at,
		    deleted_by = :deleted_by
		WHERE id = :id
`, org); err != nil {
		return fmt.Errorf("failed to update organization: %w", err)
	}
	return nil
}
