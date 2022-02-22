package db

import (
	"context"
	"fmt"

	"getsturdy.com/api/pkg/workspace"

	"github.com/jmoiron/sqlx"
)

type repo struct {
	db *sqlx.DB
}

func NewRepo(db *sqlx.DB) Repository {
	return &repo{db: db}
}

func (r *repo) Create(entity workspace.Workspace) error {
	_, err := r.db.NamedExec(`INSERT INTO workspaces
		(id, user_id, codebase_id, name, created_at, view_id, latest_snapshot_id, draft_description)
		VALUES
		(:id, :user_id, :codebase_id, :name, :created_at, :view_id, :latest_snapshot_id, :draft_description)`, &entity)
	if err != nil {
		return fmt.Errorf("failed to perform insert: %w", err)
	}
	return nil
}

func (r *repo) Get(id string) (*workspace.Workspace, error) {
	var entity workspace.Workspace
	err := r.db.Get(&entity, `SELECT id, user_id, codebase_id, name, ready_for_review_change, approved_change, created_at, last_landed_at, archived_at, unarchived_at, updated_at, draft_description, view_id, latest_snapshot_id, up_to_date_with_trunk, head_change_id, head_change_computed
	FROM workspaces
	WHERE id=$1`, id)
	if err != nil {
		return nil, fmt.Errorf("failed to query table: %w", err)
	}
	return &entity, nil
}

func (r *repo) ListByCodebaseIDs(codebaseIDs []string, includeArchived bool) ([]*workspace.Workspace, error) {
	q := `SELECT id, user_id, codebase_id, name, ready_for_review_change, approved_change, created_at, last_landed_at, archived_at, unarchived_at, updated_at, draft_description, view_id, latest_snapshot_id, up_to_date_with_trunk, head_change_id, head_change_computed
	FROM workspaces
	WHERE codebase_id IN(?)`

	if !includeArchived {
		q += "  AND archived_at IS NULL"
	}

	query, args, err := sqlx.In(q, codebaseIDs)
	if err != nil {
		return nil, fmt.Errorf("failed to create query: %w", err)
	}
	query = r.db.Rebind(query)

	var views []*workspace.Workspace
	err = r.db.Select(&views, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to query table: %w", err)
	}
	return views, nil
}

func (r *repo) ListByCodebaseIDsAndUserID(codebaseIDs []string, userID string) ([]*workspace.Workspace, error) {
	query, args, err := sqlx.In(`SELECT id, user_id, codebase_id, name, ready_for_review_change, approved_change, created_at, last_landed_at, archived_at, unarchived_at, updated_at, draft_description, view_id, latest_snapshot_id, up_to_date_with_trunk, head_change_id
	FROM workspaces
	WHERE codebase_id IN(?)
	  AND user_id = ?
	  AND archived_at IS NULL`,
		codebaseIDs,
		userID)
	if err != nil {
		return nil, fmt.Errorf("failed to create query: %w", err)
	}
	query = r.db.Rebind(query)

	var views []*workspace.Workspace
	err = r.db.Select(&views, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to query table: %w", err)
	}
	return views, nil
}

func (r *repo) Update(ctx context.Context, entity *workspace.Workspace) error {
	if _, err := r.db.NamedExecContext(ctx, `UPDATE workspaces
		SET name = :name,
			ready_for_review_change = :ready_for_review_change,
			approved_change = :approved_change,
		    last_landed_at = :last_landed_at,
		    archived_at = :archived_at,
		    unarchived_at = :unarchived_at,
			up_to_date_with_trunk = :up_to_date_with_trunk,
		    updated_at = :updated_at,
		    draft_description = :draft_description,
		    view_id = :view_id,
		    latest_snapshot_id = :latest_snapshot_id,
		    head_change_id = :head_change_id,
		    head_change_computed = :head_change_computed
		WHERE id = :id`, &entity); err != nil {
		return fmt.Errorf("failed to perform update: %w", err)
	}
	return nil
}

func (r *repo) UnsetUpToDateWithTrunkForAllInCodebase(codebaseID string) error {
	_, err := r.db.Exec("UPDATE workspaces SET up_to_date_with_trunk = NULL WHERE codebase_id = $1 AND archived_at IS NULL", codebaseID)
	if err != nil {
		return fmt.Errorf("failed to perform update: %w", err)
	}
	return nil
}

func (r *repo) GetByViewID(viewID string, includeArchived bool) (*workspace.Workspace, error) {
	var entity workspace.Workspace

	q := `SELECT id, user_id, codebase_id, name, ready_for_review_change, approved_change, created_at, last_landed_at, archived_at, unarchived_at, updated_at, draft_description, view_id, latest_snapshot_id, up_to_date_with_trunk, head_change_id, head_change_computed
		FROM workspaces
		WHERE view_id=$1`

	if !includeArchived {
		q += " AND archived_at IS NULL"
	}
	err := r.db.Get(&entity, q, viewID)
	if err != nil {
		return nil, fmt.Errorf("failed to query table: %w", err)
	}
	return &entity, nil
}

func (r *repo) GetBySnapshotID(snapshotID string) (*workspace.Workspace, error) {
	var entity workspace.Workspace
	if err := r.db.Get(&entity, `
		SELECT 
			id, 
			user_id,
			codebase_id,
			name,
			ready_for_review_change,
			approved_change,
			created_at,
			last_landed_at,
			archived_at,
			unarchived_at,
			updated_at,
			draft_description,
			view_id,
			latest_snapshot_id,
			up_to_date_with_trunk,
			head_change_id,
			head_change_computed
		FROM 
			workspaces
		WHERE
			latest_snapshot_id=$1
	`, snapshotID); err != nil {
		return nil, fmt.Errorf("failed to query table: %w", err)
	}
	return &entity, nil
}
