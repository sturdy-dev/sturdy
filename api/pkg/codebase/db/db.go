package db

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"

	"getsturdy.com/api/pkg/codebase"
)

var _ CodebaseRepository = &Repo{}

type Repo struct {
	db *sqlx.DB
}

func NewRepo(db *sqlx.DB) CodebaseRepository {
	return &Repo{db: db}
}

func (r *Repo) Create(entity codebase.Codebase) error {
	_, err := r.db.NamedExec(`INSERT INTO codebases (id, short_id, name, description, emoji, created_at, invite_code, is_ready, is_public, organization_id, calculated_head_change_id, cached_head_change_id)
		VALUES (:id, :short_id, :name, :description, :emoji, :created_at, :invite_code, :is_ready, :is_public, :organization_id, :calculated_head_change_id, :cached_head_change_id)`, &entity)
	if err != nil {
		return fmt.Errorf("failed to create codebase: %w", err)
	}
	return nil
}

func (r *Repo) Get(id string) (*codebase.Codebase, error) {
	entity := &codebase.Codebase{}
	err := r.db.Get(entity, `SELECT id, short_id, name, description, emoji, created_at, invite_code, is_ready, archived_at, is_public, organization_id, calculated_head_change_id, cached_head_change_id
		FROM codebases
		WHERE id = $1
		AND archived_at IS NULL`, id)
	if err != nil {
		return nil, fmt.Errorf("failed to query table: %w", err)
	}
	return entity, nil
}

func (r *Repo) GetAllowArchived(id string) (*codebase.Codebase, error) {
	entity := &codebase.Codebase{}
	err := r.db.Get(entity, `SELECT id, short_id, name, description, emoji, created_at, invite_code, is_ready, archived_at, is_public, organization_id, calculated_head_change_id, cached_head_change_id
		FROM codebases
		WHERE id = $1`, id)
	if err != nil {
		return nil, fmt.Errorf("failed to query table: %w", err)
	}
	return entity, nil
}

func (r *Repo) GetByInviteCode(inviteCode string) (*codebase.Codebase, error) {
	entity := &codebase.Codebase{}
	err := r.db.Get(entity, `SELECT id, short_id, name, description, emoji, created_at, invite_code, is_ready, archived_at, is_public, organization_id, calculated_head_change_id, cached_head_change_id
		FROM codebases
		WHERE invite_code = $1
	    AND archived_at IS NULL`, inviteCode)
	if err != nil {
		return nil, fmt.Errorf("failed to query table: %w", err)
	}
	return entity, nil
}

func (r *Repo) GetByShortID(shortID string) (*codebase.Codebase, error) {
	entity := &codebase.Codebase{}
	err := r.db.Get(entity, `SELECT id, short_id, name, description, emoji, created_at, invite_code, is_ready, archived_at, is_public, organization_id, calculated_head_change_id, cached_head_change_id
		FROM codebases
		WHERE short_id = $1
	    AND archived_at IS NULL`, shortID)
	if err != nil {
		return nil, fmt.Errorf("failed to query table: %w", err)
	}
	return entity, nil
}

func (r *Repo) Update(entity *codebase.Codebase) error {
	_, err := r.db.NamedExec(`UPDATE codebases
		SET name = :name,
			description = :description,
			emoji = :emoji,
		    invite_code = :invite_code,
			archived_at = :archived_at,
		    short_id = :short_id,
		    is_ready = :is_ready,
		    is_public = :is_public,
		    organization_id = :organization_id,
			calculated_head_change_id = :calculated_head_change_id,
			cached_head_change_id = :cached_head_change_id
		WHERE id = :id`, &entity)
	if err != nil {
		return fmt.Errorf("failed to perform update: %w", err)
	}
	return nil
}

func (r *Repo) ListByOrganization(ctx context.Context, organizationID string) ([]*codebase.Codebase, error) {
	var res []*codebase.Codebase
	err := r.db.SelectContext(ctx, &res, `
		SELECT id, short_id, name, description, emoji, created_at, invite_code, is_ready, archived_at, is_public, organization_id, calculated_head_change_id, cached_head_change_id
		FROM codebases
		WHERE organization_id = $1
	    AND archived_at IS NULL`, organizationID)
	if err != nil {
		return nil, fmt.Errorf("failed to list codebases by organization: %w", err)
	}
	return res, nil
}

func (r *Repo) Count(ctx context.Context) (uint64, error) {
	var res struct {
		Count uint64
	}
	if err := r.db.GetContext(ctx, &res, "SELECT count(1) as Count FROM codebases"); err != nil {
		return 0, fmt.Errorf("failed to select: %w", err)
	}
	return res.Count, nil
}
