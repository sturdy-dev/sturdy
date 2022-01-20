package db

import (
	"context"
	"fmt"

	"getsturdy.com/api/pkg/workspace/activity"

	"github.com/jmoiron/sqlx"
)

type ActivityReadsRepository interface {
	Create(context.Context, activity.WorkspaceActivityReads) error
	Update(context.Context, *activity.WorkspaceActivityReads) error
	GetByUserAndWorkspace(ctx context.Context, userID, workspaceID string) (*activity.WorkspaceActivityReads, error)
}

type activityReadsRepo struct {
	db *sqlx.DB
}

func NewActivityReadsRepo(db *sqlx.DB) ActivityReadsRepository {
	return &activityReadsRepo{db: db}
}

func (r *activityReadsRepo) Create(ctx context.Context, entity activity.WorkspaceActivityReads) error {
	_, err := r.db.NamedExecContext(ctx, `INSERT INTO workspace_activity_reads
		(id, user_id, workspace_id, last_read_created_at)
		VALUES
		(:id, :user_id, :workspace_id, :last_read_created_at)`, &entity)
	if err != nil {
		return fmt.Errorf("failed to perform insert: %w", err)
	}
	return nil
}

func (r *activityReadsRepo) Update(ctx context.Context, entity *activity.WorkspaceActivityReads) error {
	_, err := r.db.NamedExecContext(ctx, `UPDATE workspace_activity_reads
		SET last_read_created_at = :last_read_created_at
		WHERE id = :id`, &entity)
	if err != nil {
		return fmt.Errorf("failed to perform insert: %w", err)
	}
	return nil
}

func (r *activityReadsRepo) GetByUserAndWorkspace(ctx context.Context, userID, workspaceID string) (*activity.WorkspaceActivityReads, error) {
	var res activity.WorkspaceActivityReads
	if err := r.db.GetContext(ctx, &res, `SELECT id, user_id, workspace_id, last_read_created_at
		FROM workspace_activity_reads
		WHERE workspace_id = $1
		AND user_id = $2`, workspaceID, userID); err != nil {
		return nil, fmt.Errorf("failed to query table: %w", err)
	}
	return &res, nil
}
