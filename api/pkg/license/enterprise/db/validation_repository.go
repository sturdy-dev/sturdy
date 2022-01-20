package db

import (
	"context"

	"github.com/jmoiron/sqlx"

	"getsturdy.com/api/pkg/license/enterprise/license"
)

type ValidationRepository interface {
	Record(ctx context.Context, validation license.SelfHostedLicenseValidation) error
}

type validationRepository struct {
	db *sqlx.DB
}

func NewValidationRepository(db *sqlx.DB) ValidationRepository {
	return &validationRepository{db: db}
}

func (r *validationRepository) Record(ctx context.Context, validation license.SelfHostedLicenseValidation) error {
	_, err := r.db.NamedExecContext(ctx, `INSERT INTO self_hosted_license_validations (id, self_hosted_license_id, validated_at, status, reported_version, reported_booted_at, reported_user_count, reported_codebase_count, from_ip_addr)
		VALUES (:id, :self_hosted_license_id, :validated_at, :status, :reported_version, :reported_booted_at, :reported_user_count, :reported_codebase_count, :from_ip_addr) 
		`, validation)
	if err != nil {
		return err
	}
	return nil
}
