package db

import (
	"fmt"

	"github.com/jmoiron/sqlx"

	"mash/pkg/codebase"
)

var _ CodebaseRepository = &Repo{}

type Repo struct {
	db *sqlx.DB
}

func NewRepo(db *sqlx.DB) CodebaseRepository {
	return &Repo{db: db}
}

func (r *Repo) Create(entity codebase.Codebase) error {
	result, err := r.db.NamedExec(`INSERT INTO codebases (id, short_id, name, description, emoji, created_at, invite_code, is_ready, is_public)
		VALUES (:id, :short_id, :name, :description, :emoji, :created_at, :invite_code, :is_ready, :is_public)`, &entity)
	if err != nil {
		return fmt.Errorf("failed to perform insert: %w", err)
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected %w", err)
	}
	if rows != 1 {
		return fmt.Errorf("unexpected number of rows affected, expected 1, actual: %d", rows)
	}
	return nil
}

func (r *Repo) Get(id string) (*codebase.Codebase, error) {
	entity := &codebase.Codebase{}
	err := r.db.Get(entity, `SELECT id, short_id, name, description, emoji, created_at, invite_code, is_ready, archived_at, is_public
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
	err := r.db.Get(entity, `SELECT id, short_id, name, description, emoji, created_at, invite_code, is_ready, archived_at, is_public
		FROM codebases
		WHERE id = $1`, id)
	if err != nil {
		return nil, fmt.Errorf("failed to query table: %w", err)
	}
	return entity, nil
}

func (r *Repo) GetByInviteCode(inviteCode string) (*codebase.Codebase, error) {
	entity := &codebase.Codebase{}
	err := r.db.Get(entity, `SELECT id, short_id, name, description, emoji, created_at, invite_code, is_ready, archived_at, is_public
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
	err := r.db.Get(entity, `SELECT id, short_id, name, description, emoji, created_at, invite_code, is_ready, archived_at, is_public
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
		    is_public = :is_public
		WHERE id = :id`, &entity)
	if err != nil {
		return fmt.Errorf("failed to perform update: %w", err)
	}
	return nil
}
