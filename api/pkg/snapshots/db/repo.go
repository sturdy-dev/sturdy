package db

import (
	"context"
	"fmt"

	"getsturdy.com/api/pkg/snapshots"

	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
)

type Repository interface {
	Create(snapshot *snapshots.Snapshot) error
	LatestInWorkspace(context.Context, string) (*snapshots.Snapshot, error)
	GetByCommitSHA(context.Context, string) (*snapshots.Snapshot, error)
	ListByIDs(context.Context, []string) ([]*snapshots.Snapshot, error)
	Get(string) (*snapshots.Snapshot, error)
	Update(snapshot *snapshots.Snapshot) error
}

type dbrepo struct {
	db *sqlx.DB
}

func NewRepo(db *sqlx.DB) Repository {
	return &dbrepo{db: db}
}

func (r *dbrepo) Create(snapshot *snapshots.Snapshot) error {
	_, err := r.db.NamedExec(`INSERT INTO snapshots
		(id, created_at, previous_snapshot_id, codebase_id, commit_id, workspace_id , action, diffs_count)
		VALUES(:id, :created_at, :previous_snapshot_id, :codebase_id, :commit_id, :workspace_id, :action, :diffs_count)
    	`, &snapshot)
	if err != nil {
		return fmt.Errorf("failed to insert: %w", err)
	}
	return nil
}

func (r *dbrepo) Get(id string) (*snapshots.Snapshot, error) {
	var res snapshots.Snapshot
	err := r.db.Get(&res, `SELECT id, created_at, previous_snapshot_id, codebase_id, commit_id, workspace_id,  action, diffs_count
		FROM snapshots
		WHERE id=$1
		AND deleted_at IS NULL`, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get: %w", err)
	}
	return &res, nil
}

func (r *dbrepo) LatestInWorkspace(ctx context.Context, workspaceID string) (*snapshots.Snapshot, error) {
	var res snapshots.Snapshot
	if err := r.db.GetContext(ctx, &res, `SELECT id, created_at, previous_snapshot_id, codebase_id, commit_id, workspace_id,  action, diffs_count
		FROM snapshots
		WHERE workspace_id = $1
		AND deleted_at IS NULL
		ORDER BY created_at DESC
		LIMIT 1
	`, workspaceID); err != nil {
		return nil, fmt.Errorf("failed to get latest in workspace: %w", err)
	}
	return &res, nil
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

// GetByCommitSHA returns snapshots by commit_id
func (r *dbrepo) GetByCommitSHA(ctx context.Context, sha string) (*snapshots.Snapshot, error) {
	var res snapshots.Snapshot
	if err := r.db.GetContext(ctx, &res, `SELECT id, created_at, previous_snapshot_id, codebase_id, commit_id, workspace_id,  action, diffs_count
		FROM snapshots
		WHERE commit_id=$1`, sha); err != nil {
		return nil, fmt.Errorf("failed to get by commit sha: %w", err)
	}
	return &res, nil
}

func (r *dbrepo) ListByIDs(ctx context.Context, ids []string) ([]*snapshots.Snapshot, error) {
	if len(ids) == 0 {
		return nil, nil
	}
	var res []*snapshots.Snapshot
	if err := r.db.SelectContext(ctx, &res, `SELECT id, created_at, previous_snapshot_id, codebase_id, commit_id, workspace_id,  action, diffs_count
		FROM snapshots
		WHERE id = ANY($1)`, pq.Array(ids)); err != nil {
		return nil, fmt.Errorf("failed to list by ids: %w", err)
	}
	return res, nil
}
