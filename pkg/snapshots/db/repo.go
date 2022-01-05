package db

import (
	"fmt"
	"mash/pkg/snapshots"

	"github.com/jmoiron/sqlx"
)

type Repository interface {
	Create(snapshot *snapshots.Snapshot) error
	ListByView(viewID string) ([]*snapshots.Snapshot, error)
	LatestInView(viewID string) (*snapshots.Snapshot, error)
	LatestInViewAndWorkspace(viewID, workspaceID string) (*snapshots.Snapshot, error)
	Get(ID string) (*snapshots.Snapshot, error)
	Update(snapshot *snapshots.Snapshot) error
	ListUndeletedInCodebase(codebaseID string) ([]*snapshots.Snapshot, error)
}

type dbrepo struct {
	db *sqlx.DB
}

func NewRepo(db *sqlx.DB) Repository {
	return &dbrepo{db: db}
}

func (r *dbrepo) Create(snapshot *snapshots.Snapshot) error {
	_, err := r.db.NamedExec(`INSERT INTO snapshots
		(id, view_id, created_at, new_files, changed_files, deleted_files, previous_snapshot_id, codebase_id, commit_id, workspace_id , action)
		VALUES(:id, :view_id, :created_at, :new_files, :changed_files, :deleted_files, :previous_snapshot_id, :codebase_id, :commit_id, :workspace_id, :action)
    	`, &snapshot)
	if err != nil {
		return fmt.Errorf("failed to insert: %w", err)
	}
	return nil
}

func (r *dbrepo) Get(id string) (*snapshots.Snapshot, error) {
	var res snapshots.Snapshot
	err := r.db.Get(&res, `SELECT id, view_id, created_at, new_files, changed_files, deleted_files, previous_snapshot_id, codebase_id, commit_id, workspace_id,  action
		FROM snapshots
		WHERE id=$1
		AND deleted_at IS NULL`, id)
	if err != nil {
		return nil, fmt.Errorf("failed to query table: %w", err)
	}
	return &res, nil
}

func (r *dbrepo) LatestInView(viewID string) (*snapshots.Snapshot, error) {
	var res snapshots.Snapshot
	err := r.db.Get(&res, `SELECT id, view_id, created_at, new_files, changed_files, deleted_files, previous_snapshot_id, codebase_id, commit_id, workspace_id,  action
		FROM snapshots
		WHERE view_id=$1
		AND deleted_at IS NULL
		ORDER BY created_at DESC
		LIMIT 1`, viewID)
	if err != nil {
		return nil, fmt.Errorf("failed to query table: %w", err)
	}
	return &res, nil
}

func (r *dbrepo) LatestInViewAndWorkspace(viewID, workspaceID string) (*snapshots.Snapshot, error) {
	var res snapshots.Snapshot
	err := r.db.Get(&res, `SELECT id, view_id, created_at, new_files, changed_files, deleted_files, previous_snapshot_id, codebase_id, commit_id, workspace_id,  action
		FROM snapshots
		WHERE view_id=$1
	    AND workspace_id=$2
		AND deleted_at IS NULL
		ORDER BY created_at DESC
		LIMIT 1`, viewID, workspaceID)
	if err != nil {
		return nil, fmt.Errorf("failed to query table: %w", err)
	}
	return &res, nil
}

func (r *dbrepo) ListByView(viewID string) ([]*snapshots.Snapshot, error) {
	var res []*snapshots.Snapshot
	err := r.db.Select(&res, `SELECT id, view_id, created_at, new_files, changed_files, deleted_files, previous_snapshot_id, codebase_id, commit_id, workspace_id,  action
		FROM snapshots
		WHERE view_id=$1
		  AND deleted_at IS NULL
		ORDER BY created_at DESC
		LIMIT 100`, viewID)
	if err != nil {
		return nil, fmt.Errorf("failed to query table: %w", err)
	}
	return res, nil
}

func (r *dbrepo) ListUndeletedInCodebase(codebaseID string) ([]*snapshots.Snapshot, error) {
	var res []*snapshots.Snapshot
	if err := r.db.Select(&res, `
		SELECT 
			id,
			view_id,
			created_at,
			new_files,
			changed_files,
			deleted_files,
			previous_snapshot_id,
			codebase_id,
			commit_id,
			workspace_id,
			action
		FROM 
			snapshots
		WHERE codebase_id=$1
	      AND deleted_at IS NULL
		ORDER BY 
		  created_at DESC
		LIMIT 50000
		`, codebaseID); err != nil {
		return nil, fmt.Errorf("failed to query table: %w", err)
	}
	return res, nil
}

func (r *dbrepo) Update(snapshot *snapshots.Snapshot) error {
	_, err := r.db.NamedExec(`UPDATE snapshots
		SET deleted_at = :deleted_at
		WHERE id = :id
    	`, &snapshot)
	if err != nil {
		return fmt.Errorf("failed to update: %w", err)
	}
	return nil
}
