package db

import (
	"context"
	"fmt"
	"mash/pkg/codebase/acl"

	"github.com/jmoiron/sqlx"
)

type ACLRepository interface {
	Create(context.Context, acl.ACL) error
	Update(context.Context, acl.ACL) error
	GetByCodebaseID(ctx context.Context, codebaseID string) (acl.ACL, error)
}

type aclRepository struct {
	db *sqlx.DB
}

func NewACLRepository(db *sqlx.DB) ACLRepository {
	return &aclRepository{
		db: db,
	}
}

func (r *aclRepository) Create(ctx context.Context, entity acl.ACL) error {
	result, err := r.db.NamedExecContext(ctx, `INSERT INTO acls
		(id, codebase_id, created_at, policy)
		VALUES
		(:id, :codebase_id, :created_at, :policy)`, entity)
	if err != nil {
		return fmt.Errorf("failed to perform insert: %w", err)
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	if rows != 1 {
		return fmt.Errorf("unexpected number of rows affected, expected 1, got %d", rows)
	}
	return nil
}

func (r *aclRepository) Update(ctx context.Context, entity acl.ACL) error {
	result, err := r.db.NamedExecContext(ctx, `UPDATE acls
		SET policy = :policy
		WHERE id = :id`, &entity)
	if err != nil {
		return fmt.Errorf("failed to perform update: %w", err)
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	if rows != 1 {
		return fmt.Errorf("unexpected number of rows affected, expected 1, got %d", rows)
	}
	return nil
}

func (r *aclRepository) GetByCodebaseID(ctx context.Context, codebaseID string) (acl.ACL, error) {
	entity := &acl.ACL{}
	err := r.db.GetContext(ctx, entity, `SELECT id, codebase_id, created_at, policy
		FROM acls
		WHERE codebase_id = $1`, codebaseID)
	if err != nil {
		return acl.ACL{}, fmt.Errorf("failed to query table: %w", err)
	}
	return *entity, nil
}
