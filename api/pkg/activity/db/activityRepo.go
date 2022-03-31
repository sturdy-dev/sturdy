package db

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"getsturdy.com/api/pkg/activity"
	"getsturdy.com/api/pkg/changes"

	"github.com/jmoiron/sqlx"
)

type ActivityRepository interface {
	Create(context.Context, activity.Activity) error
	Get(ctx context.Context, id string) (*activity.Activity, error)
	ListByWorkspaceID(ctx context.Context, workspaceID string, limit int32) ([]*activity.Activity, error)
	ListByWorkspaceIDNewerThan(ctx context.Context, workspaceID string, newerThan time.Time, limit int32) ([]*activity.Activity, error)

	SetChangeID(ctx context.Context, workspaceID string, changeID changes.ID) error
	ListByChangeID(context.Context, changes.ID, int32) ([]*activity.Activity, error)
}

type activityRepo struct {
	db *sqlx.DB
}

func NewActivityRepository(db *sqlx.DB) ActivityRepository {
	return &activityRepo{db: db}
}

func (r *activityRepo) Create(ctx context.Context, entity activity.Activity) error {
	_, err := r.db.NamedExecContext(ctx, `INSERT INTO workspace_activity
		(id, user_id, workspace_id, created_at, activity_type, reference, change_id)
		VALUES
		(:id, :user_id, :workspace_id, :created_at, :activity_type, :reference, :change_id)`, &entity)
	if err != nil {
		return fmt.Errorf("failed to perform insert: %w", err)
	}
	return nil
}

func (r *activityRepo) Get(ctx context.Context, id string) (*activity.Activity, error) {
	var res activity.Activity
	if err := r.db.GetContext(ctx, &res, `SELECT id, user_id, workspace_id, created_at, activity_type, reference, change_id
		FROM workspace_activity
		WHERE id = $1`, id); err != nil {
		return nil, fmt.Errorf("failed to query table: %w", err)
	}
	return &res, nil
}

func (r *activityRepo) ListByWorkspaceID(ctx context.Context, workspaceID string, limit int32) ([]*activity.Activity, error) {
	var activities []*activity.Activity
	if err := r.db.SelectContext(ctx, &activities, `SELECT id, user_id, workspace_id, created_at, activity_type, reference, change_id
		FROM workspace_activity
		WHERE workspace_id = $1
		ORDER BY created_at DESC
		LIMIT $2`, workspaceID, limit); err != nil {
		return nil, fmt.Errorf("failed to query table: %w", err)
	}
	return activities, nil
}

func (r *activityRepo) ListByWorkspaceIDNewerThan(ctx context.Context, workspaceID string, newerThan time.Time, limit int32) ([]*activity.Activity, error) {
	var activities []*activity.Activity
	if err := r.db.SelectContext(ctx, &activities, `SELECT id, user_id, workspace_id, created_at, activity_type, reference, change_id
		FROM workspace_activity
		WHERE workspace_id = $1
		AND created_at > $2
		ORDER BY created_at DESC
		LIMIT $3`, workspaceID, newerThan, limit); err != nil {
		return nil, fmt.Errorf("failed to query table: %w", err)
	}
	return activities, nil
}

func (r *activityRepo) SetChangeID(ctx context.Context, workspaceID string, changeID changes.ID) error {
	if _, err := r.db.ExecContext(ctx, `UPDATE workspace_activity
		SET change_id = $1
		WHERE workspace_id = $2 AND change_id IS NULL
	`, changeID, workspaceID); err != nil {
		return fmt.Errorf("failed to perform update: %w", err)
	}
	return nil
}

func (r *activityRepo) ListByChangeID(ctx context.Context, changeID changes.ID, limit int32) ([]*activity.Activity, error) {
	var activities []*activity.Activity
	if err := r.db.SelectContext(ctx, &activities, `SELECT id, user_id, workspace_id, created_at, activity_type, reference, change_id
		FROM workspace_activity
		WHERE change_id = $1
		ORDER BY created_at DESC
		LIMIT $2`, changeID, limit); err != nil {
		return nil, fmt.Errorf("failed to query table: %w", err)
	}
	return activities, nil
}

type inmemory struct {
	byID          map[string]*activity.Activity
	byWorkspaceID map[string][]*activity.Activity
	byChangeID    map[changes.ID][]*activity.Activity
}

func NewInMemoryRepo() ActivityRepository {
	return &inmemory{
		byID:          make(map[string]*activity.Activity),
		byWorkspaceID: make(map[string][]*activity.Activity),
		byChangeID:    make(map[changes.ID][]*activity.Activity),
	}
}

func (i *inmemory) Create(ctx context.Context, activity activity.Activity) error {
	i.byID[activity.ID] = &activity
	if activity.WorkspaceID != nil {
		i.byWorkspaceID[*activity.WorkspaceID] = append(i.byWorkspaceID[*activity.WorkspaceID], &activity)
	}
	if activity.ChangeID != nil {
		i.byChangeID[*activity.ChangeID] = append(i.byChangeID[*activity.ChangeID], &activity)
	}
	return nil
}

func (i *inmemory) Get(ctx context.Context, id string) (*activity.Activity, error) {
	activity, ok := i.byID[id]
	if !ok {
		return nil, sql.ErrNoRows
	}
	return activity, nil
}

func (i *inmemory) ListByWorkspaceID(ctx context.Context, workspaceID string, limit int32) ([]*activity.Activity, error) {
	activities := i.byWorkspaceID[workspaceID]
	if len(activities) > int(limit) {
		return activities[:limit], nil
	}
	return activities, nil
}

func (i *inmemory) ListByWorkspaceIDNewerThan(ctx context.Context, workspaceID string, newerThan time.Time, limit int32) ([]*activity.Activity, error) {
	var activities []*activity.Activity
	for _, activity := range i.byWorkspaceID[workspaceID] {
		if activity.CreatedAt.After(newerThan) {
			activities = append(activities, activity)
		}
	}
	if len(activities) > int(limit) {
		return activities[:limit], nil
	}
	return activities, nil
}

func (i *inmemory) SetChangeID(ctx context.Context, workspaceID string, changeID changes.ID) error {
	for _, activity := range i.byWorkspaceID[workspaceID] {
		if activity.ChangeID == nil {
			activity.ChangeID = &changeID
			i.byChangeID[changeID] = append(i.byChangeID[changeID], activity)
		}
	}
	return nil
}

func (i *inmemory) ListByChangeID(ctx context.Context, changeID changes.ID, limit int32) ([]*activity.Activity, error) {
	activities := i.byChangeID[changeID]
	if len(activities) > int(limit) {
		return activities[:limit], nil
	}
	return activities, nil
}
