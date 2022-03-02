package db

import (
	"context"
	"fmt"
	"time"

	"getsturdy.com/api/pkg/activity"

	"github.com/jmoiron/sqlx"
)

type ActivityRepository interface {
	Create(context.Context, activity.Activity) error
	Get(ctx context.Context, id string) (*activity.Activity, error)
	ListByWorkspaceID(ctx context.Context, workspaceID string) ([]*activity.Activity, error)
	ListByWorkspaceIDNewerThan(ctx context.Context, workspaceID string, newerThan time.Time) ([]*activity.Activity, error)
}

type activityRepo struct {
	db *sqlx.DB
}

func NewActivityRepo(db *sqlx.DB) ActivityRepository {
	return &activityRepo{db: db}
}

func (r *activityRepo) Create(ctx context.Context, entity activity.Activity) error {
	_, err := r.db.NamedExecContext(ctx, `INSERT INTO workspace_activity
		(id, user_id, workspace_id, created_at, activity_type, reference)
		VALUES
		(:id, :user_id, :workspace_id, :created_at, :activity_type, :reference)`, &entity)
	if err != nil {
		return fmt.Errorf("failed to perform insert: %w", err)
	}
	return nil
}

func (r *activityRepo) Get(ctx context.Context, id string) (*activity.Activity, error) {
	var res activity.Activity
	if err := r.db.GetContext(ctx, &res, `SELECT id, user_id, workspace_id, created_at, activity_type, reference
		FROM workspace_activity
		WHERE id = $1`, id); err != nil {
		return nil, fmt.Errorf("failed to query table: %w", err)
	}
	return &res, nil
}

func (r *activityRepo) ListByWorkspaceID(ctx context.Context, workspaceID string) ([]*activity.Activity, error) {
	var activities []*activity.Activity
	if err := r.db.SelectContext(ctx, &activities, `SELECT id, user_id, workspace_id, created_at, activity_type, reference
		FROM workspace_activity
		WHERE workspace_id = $1
		ORDER BY created_at DESC`, workspaceID); err != nil {
		return nil, fmt.Errorf("failed to query table: %w", err)
	}
	return activities, nil
}

func (r *activityRepo) ListByWorkspaceIDNewerThan(ctx context.Context, workspaceID string, newerThan time.Time) ([]*activity.Activity, error) {
	var activities []*activity.Activity
	if err := r.db.SelectContext(ctx, &activities, `SELECT id, user_id, workspace_id, created_at, activity_type, reference
		FROM workspace_activity
		WHERE workspace_id = $1
		AND created_at > $2
		ORDER BY created_at DESC`, workspaceID, newerThan); err != nil {
		return nil, fmt.Errorf("failed to query table: %w", err)
	}
	return activities, nil
}
