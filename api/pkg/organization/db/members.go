package db

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"

	"getsturdy.com/api/pkg/organization"
	"getsturdy.com/api/pkg/users"
)

type MemberRepository interface {
	GetByUserIDAndOrganizationID(ctx context.Context, userID users.ID, organizationID string) (*organization.Member, error)
	ListByOrganizationID(ctx context.Context, id string) ([]*organization.Member, error)
	ListByUserID(context.Context, users.ID) ([]*organization.Member, error)
	Create(ctx context.Context, org organization.Member) error
	Update(ctx context.Context, org *organization.Member) error
	GetByID(ctx context.Context, id string) (*organization.Member, error)
}

type memberRepository struct {
	db *sqlx.DB
}

func NewMember(db *sqlx.DB) MemberRepository {
	return &memberRepository{db: db}
}

func (r *memberRepository) GetByID(ctx context.Context, id string) (*organization.Member, error) {
	var mem organization.Member
	if err := r.db.GetContext(ctx, &mem, `SELECT id, user_id, organization_id, created_at, created_by, deleted_at, deleted_by
		FROM organization_members
		WHERE id = $1
		  AND deleted_at IS NULL`, id); err != nil {
		return nil, fmt.Errorf("failed to get organization_member by id: %w", err)
	}
	return &mem, nil
}

func (r *memberRepository) GetByUserIDAndOrganizationID(ctx context.Context, userID users.ID, organizationID string) (*organization.Member, error) {
	var mem organization.Member
	if err := r.db.GetContext(ctx, &mem, `SELECT id, user_id, organization_id, created_at, created_by, deleted_at, deleted_by
		FROM organization_members
		WHERE user_id = $1
		  AND organization_id = $2
		  AND deleted_at IS NULL`, userID, organizationID); err != nil {
		return nil, fmt.Errorf("could not get organization_members: %w", err)
	}
	return &mem, nil
}

func (r *memberRepository) ListByOrganizationID(ctx context.Context, id string) ([]*organization.Member, error) {
	var res []*organization.Member
	if err := r.db.SelectContext(ctx, &res, `SELECT id, user_id, organization_id, created_at, created_by, deleted_at, deleted_by
		FROM organization_members
		WHERE organization_id = $1
		  AND deleted_at IS NULL`, id); err != nil {
		return nil, fmt.Errorf("failed to list organization_members by org_id: %w", err)
	}
	return res, nil
}

func (r *memberRepository) ListByUserID(ctx context.Context, id users.ID) ([]*organization.Member, error) {
	var res []*organization.Member
	if err := r.db.SelectContext(ctx, &res, `SELECT id, user_id, organization_id, created_at, created_by, deleted_at, deleted_by
		FROM organization_members
		WHERE user_id = $1
		  AND deleted_at IS NULL`, id); err != nil {
		return nil, fmt.Errorf("failed to list organization_members by user_id: %w", err)
	}
	return res, nil
}

func (r *memberRepository) Create(ctx context.Context, mem organization.Member) error {
	if _, err := r.db.NamedExecContext(ctx, `INSERT INTO organization_members (id, user_id, organization_id, created_at, created_by, deleted_at, deleted_by)
		VALUES (:id, :user_id, :organization_id, :created_at, :created_by, :deleted_at, :deleted_by)
		ON CONFLICT (organization_id, user_id) DO UPDATE
		SET deleted_at = NULL,
		    deleted_by = NULL`,
		mem); err != nil {
		return fmt.Errorf("failed to create organization: %w", err)
	}
	return nil
}

func (r *memberRepository) Update(ctx context.Context, org *organization.Member) error {
	if _, err := r.db.NamedExecContext(ctx, `UPDATE organization_members
		SET deleted_at = :deleted_at,
		    deleted_by = :deleted_by
		WHERE id = :id
`, org); err != nil {
		return fmt.Errorf("failed to update organization: %w", err)
	}
	return nil
}
