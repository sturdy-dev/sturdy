package db

import (
	"context"
	"fmt"

	"getsturdy.com/api/pkg/changes"
	"getsturdy.com/api/pkg/workspaces"

	"github.com/jmoiron/sqlx"
)

type repo struct {
	db *sqlx.DB
}

func NewRepo(db *sqlx.DB) Repository {
	return &repo{db: db}
}

func (r *repo) Create(entity workspaces.Workspace) error {
	_, err := r.db.NamedExec(`INSERT INTO workspaces
		(id, user_id, codebase_id, name, created_at, view_id, latest_snapshot_id, draft_description, diffs_count)
		VALUES
		(:id, :user_id, :codebase_id, :name, :created_at, :view_id, :latest_snapshot_id, :draft_description, :diffs_count)`, &entity)
	if err != nil {
		return fmt.Errorf("failed to perform insert: %w", err)
	}
	return nil
}

func (r *repo) Get(id string) (*workspaces.Workspace, error) {
	var entity workspaces.Workspace
	err := r.db.Get(&entity, `SELECT id, user_id, codebase_id, name,  created_at, last_landed_at, archived_at, unarchived_at, updated_at, draft_description, view_id, latest_snapshot_id, up_to_date_with_trunk, head_change_id, head_change_computed, diffs_count, change_id
	FROM workspaces
	WHERE id=$1`, id)
	if err != nil {
		return nil, fmt.Errorf("failed to query table: %w", err)
	}
	return &entity, nil
}

func (r *repo) ListByCodebaseIDs(codebaseIDs []string, includeArchived bool) ([]*workspaces.Workspace, error) {
	q := `SELECT id, user_id, codebase_id, name, created_at, last_landed_at, archived_at, unarchived_at, updated_at, draft_description, view_id, latest_snapshot_id, up_to_date_with_trunk, head_change_id, head_change_computed, diffs_count, change_id
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

	var views []*workspaces.Workspace
	err = r.db.Select(&views, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to query table: %w", err)
	}
	return views, nil
}

func (r *repo) ListByCodebaseIDsAndUserID(codebaseIDs []string, userID string) ([]*workspaces.Workspace, error) {
	query, args, err := sqlx.In(`SELECT id, user_id, codebase_id, name, created_at, last_landed_at, archived_at, unarchived_at, updated_at, draft_description, view_id, latest_snapshot_id, up_to_date_with_trunk, head_change_id, diffs_count, change_id
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

	var views []*workspaces.Workspace
	err = r.db.Select(&views, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to query table: %w", err)
	}
	return views, nil
}

func (r *repo) SetUpToDateWithTrunk(ctx context.Context, workspaceID string, upToDateWithTrunk bool) error {
	if _, err := r.db.ExecContext(ctx, `UPDATE workspaces
		SET up_to_date_with_trunk = $1
		WHERE id = $2`, upToDateWithTrunk, workspaceID); err != nil {
		return fmt.Errorf("failed to perform update: %w", err)
	}
	return nil
}

func (r *repo) Update(ctx context.Context, entity *workspaces.Workspace) error {
	if _, err := r.db.NamedExecContext(ctx, `UPDATE workspaces
		SET name = :name,
		    last_landed_at = :last_landed_at,
		    archived_at = :archived_at,
		    unarchived_at = :unarchived_at,
			up_to_date_with_trunk = :up_to_date_with_trunk,
		    updated_at = :updated_at,
		    draft_description = :draft_description,
		    view_id = :view_id,
		    latest_snapshot_id = :latest_snapshot_id,
		    head_change_id = :head_change_id,
		    head_change_computed = :head_change_computed,
			diffs_count = :diffs_count,
			change_id = :change_id
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

func (r *repo) GetByViewID(viewID string, includeArchived bool) (*workspaces.Workspace, error) {
	var entity workspaces.Workspace

	q := `SELECT id, user_id, codebase_id, name, created_at, last_landed_at, archived_at, unarchived_at, updated_at, draft_description, view_id, latest_snapshot_id, up_to_date_with_trunk, head_change_id, head_change_computed, diffs_count, change_id
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

func (r *repo) GetBySnapshotID(snapshotID string) (*workspaces.Workspace, error) {
	var entity workspaces.Workspace
	if err := r.db.Get(&entity, `
		SELECT 
			id, 
			user_id,
			codebase_id,
			name,
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
			head_change_computed,
			diffs_count,
			change_id
		FROM 
			workspaces
		WHERE
			latest_snapshot_id=$1
	`, snapshotID); err != nil {
		return nil, fmt.Errorf("failed to query table: %w", err)
	}
	return &entity, nil
}

func (r *repo) SetHeadChange(ctx context.Context, workspaceID string, changeID *changes.ID) error {
	if _, err := r.db.ExecContext(ctx, `UPDATE workspaces
		SET head_change_id = $1,
			head_change_computed = TRUE
		WHERE id = $2`, changeID, workspaceID); err != nil {
		return fmt.Errorf("failed to perform update: %w", err)
	}
	return nil
}
