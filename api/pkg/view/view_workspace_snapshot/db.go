package view_workspace_snapshot

import (
	"fmt"
	"getsturdy.com/api/pkg/view"

	"github.com/jmoiron/sqlx"
)

type Repository interface {
	Create(snapshot view.ViewWorkspaceSnapshot) error
	Get(ViewID string, WorkspaceID string) (*view.ViewWorkspaceSnapshot, error)
	Update(snap *view.ViewWorkspaceSnapshot) error
}

type repo struct {
	db *sqlx.DB
}

func NewRepo(db *sqlx.DB) Repository {
	return &repo{db: db}
}

func (r *repo) Create(snap view.ViewWorkspaceSnapshot) error {
	_, err := r.db.NamedExec(`INSERT INTO view_workspace_snapshots (id, view_id, workspace_id, snapshot_id, created_at, updated_at)
		VALUES (:id, :view_id, :workspace_id, :snapshot_id, :created_at, :updated_at)`, &snap)
	if err != nil {
		return fmt.Errorf("failed to perform insert: %w", err)
	}
	return nil
}

func (r *repo) Get(viewID string, workspaceID string) (*view.ViewWorkspaceSnapshot, error) {
	var res view.ViewWorkspaceSnapshot
	err := r.db.Get(&res, `SELECT id, view_id, workspace_id, snapshot_id, created_at, updated_at
		FROM view_workspace_snapshots
		WHERE view_id=$1 AND workspace_id = $2
		ORDER BY created_at DESC
		LIMIT 1`, viewID, workspaceID)
	if err != nil {
		return nil, fmt.Errorf("failed to query table: %w", err)
	}
	return &res, nil
}

func (r *repo) Update(snap *view.ViewWorkspaceSnapshot) error {
	_, err := r.db.NamedExec(`UPDATE view_workspace_snapshots
		SET snapshot_id = :snapshot_id,
		    updated_at = :updated_at
		WHERE id = :id`, snap)
	if err != nil {
		return fmt.Errorf("failed to update %w", err)
	}
	return nil
}
