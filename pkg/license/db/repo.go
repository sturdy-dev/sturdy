package db

import (
	"context"

	"github.com/jmoiron/sqlx"

	"mash/pkg/license"
)

type Repository interface {
	Get(ctx context.Context, id string) (*license.SelfHostedLicense, error)
	Create(ctx context.Context, license license.SelfHostedLicense) error
	ListByCloudOrganizationID(ctx context.Context, cloudOrganizationID string) ([]*license.SelfHostedLicense, error)
}

type repository struct {
	db *sqlx.DB
}

func New(db *sqlx.DB) Repository {
	return &repository{db: db}
}

func (r *repository) Create(ctx context.Context, license license.SelfHostedLicense) error {
	_, err := r.db.ExecContext(ctx, `INSERT INTO self_hosted_licenses (id, cloud_organization_id, seats, created_at, active)
VALUES (:id, :cloud_organization_id, :seats, :created_at, :active)`, license)
	if err != nil {
		return err
	}
	return nil
}

func (r *repository) Get(ctx context.Context, id string) (*license.SelfHostedLicense, error) {
	var res license.SelfHostedLicense
	err := r.db.GetContext(ctx, &res, `SELECT id, cloud_organization_id, seats, created_at, active
		FROM self_hosted_licenses
		WHERE id = $1`, id)
	if err != nil {
		return nil, err
	}
	return &res, nil
}

func (r *repository) ListByCloudOrganizationID(ctx context.Context, cloudOrganizationID string) ([]*license.SelfHostedLicense, error) {
	var res []*license.SelfHostedLicense
	err := r.db.SelectContext(ctx, &res, `SELECT id, cloud_organization_id, seats, created_at, active
		FROM self_hosted_licenses
		WHERE cloud_organization_id = $1`, cloudOrganizationID)
	if err != nil {
		return nil, err
	}
	return res, nil
}
