package db

import (
	"context"
	"fmt"

	"getsturdy.com/api/pkg/workspaces/watchers"

	"github.com/jmoiron/sqlx"
)

var _ Repository = &database{}

type database struct {
	db *sqlx.DB
}

func NewDB(db *sqlx.DB) Repository {
	return &database{db: db}
}

func (d *database) Create(ctx context.Context, watcher *watchers.Watcher) error {
	if _, err := d.db.NamedExecContext(ctx, `
		INSERT INTO workspace_watchers 
			(workspace_id, user_id, status, created_at)
		VALUES
			(:workspace_id, :user_id, :status, :created_at)
	`, watcher); err != nil {
		return fmt.Errorf("failed to insert: %w", err)
	}
	return nil
}

func (d *database) ListWatchingByWorkspaceID(ctx context.Context, workspaceID string) ([]*watchers.Watcher, error) {
	var watchers []*watchers.Watcher
	if err := d.db.SelectContext(ctx, &watchers, `
		WITH latest AS (
			SELECT 
				workspace_id, user_id, MAX(created_at) AS created_at
			FROM
				workspace_watchers
			WHERE
				workspace_id = $1
			GROUP BY	
				workspace_id, user_id
		)

		SELECT
			workspace_watchers.workspace_id,
			workspace_watchers.user_id,
			workspace_watchers.status,
			workspace_watchers.created_at
		FROM 
			workspace_watchers JOIN latest ON
					workspace_watchers.workspace_id   = latest.workspace_id
					AND workspace_watchers.user_id    = latest.user_id
					AND workspace_watchers.created_at = latest.created_at
		WHERE
			status = 'watching'
	`, workspaceID); err != nil {
		return nil, fmt.Errorf("failed to select: %w", err)
	}
	return watchers, nil
}

func (d *database) GetByUserIDAndWorkspaceID(ctx context.Context, userID string, workspaceID string) (*watchers.Watcher, error) {
	watcher := &watchers.Watcher{}
	if err := d.db.GetContext(ctx, watcher, `
		SELECT
			workspace_id,
			user_id,
			status,
			created_at
		FROM
			workspace_watchers
		WHERE
			workspace_id = $1
			AND user_id = $2
		`, workspaceID, userID); err != nil {
		return nil, fmt.Errorf("failed to select: %w", err)
	}
	return watcher, nil
}
