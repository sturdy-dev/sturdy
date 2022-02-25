package db

import (
	"context"
	"fmt"

	"getsturdy.com/api/pkg/users"
	"getsturdy.com/api/pkg/view"

	"github.com/jmoiron/sqlx"
)

type Repository interface {
	Create(entity view.View) error
	Get(id string) (*view.View, error)
	ListByCodebase(codebaseID string) ([]*view.View, error)
	ListByUser(users.ID) ([]*view.View, error)
	LastUsedByCodebaseAndUser(ctx context.Context, codebaseID string, userID users.ID) (*view.View, error)
	ListByCodebaseAndUser(codebaseID string, userID users.ID) ([]*view.View, error)
	ListByCodebaseAndWorkspace(codebaseID, workspaceID string) ([]*view.View, error)
	Update(e *view.View) error
}

type repo struct {
	db *sqlx.DB
}

func NewRepo(db *sqlx.DB) Repository {
	return &repo{db: db}
}

func (r *repo) Create(entity view.View) error {
	result, err := r.db.NamedExec(`INSERT INTO views
    	(id, user_id, codebase_id, workspace_id, name, last_used_at, created_at, mount_path, mount_hostname)
    	VALUES(:id, :user_id, :codebase_id, :workspace_id, :name, :last_used_at, :created_at, :mount_path, :mount_hostname)
    	`, &entity)
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

func (r *repo) Get(id string) (*view.View, error) {
	var entity view.View
	err := r.db.Get(&entity, "SELECT id, user_id, codebase_id, workspace_id, name, last_used_at, created_at, mount_path, mount_hostname   FROM views WHERE id=$1", id)
	if err != nil {
		return nil, fmt.Errorf("failed to query table: %w", err)
	}
	return &entity, nil
}

func (r *repo) ListByCodebase(codebaseID string) ([]*view.View, error) {
	var views []*view.View
	err := r.db.Select(&views, "SELECT id, user_id, codebase_id, workspace_id, name, last_used_at, created_at, mount_path, mount_hostname   FROM views WHERE codebase_id=$1", codebaseID)
	if err != nil {
		return nil, fmt.Errorf("failed to query table: %w", err)
	}
	return views, nil
}

func (r *repo) ListByUser(userID users.ID) ([]*view.View, error) {
	var views []*view.View
	err := r.db.Select(&views, "SELECT id, user_id, codebase_id, workspace_id, name, last_used_at, created_at, mount_path, mount_hostname  FROM views WHERE user_id=$1", userID)
	if err != nil {
		return nil, fmt.Errorf("failed to query table: %w", err)
	}
	return views, nil
}

func (r *repo) LastUsedByCodebaseAndUser(ctx context.Context, codebaseID string, userID users.ID) (*view.View, error) {
	var entity view.View
	err := r.db.GetContext(ctx, &entity, "SELECT id, user_id, codebase_id, workspace_id, name, last_used_at, created_at, mount_path, mount_hostname FROM views WHERE user_id=$1 AND codebase_id=$2 ORDER BY last_used_at DESC LIMIT 1", userID, codebaseID)
	if err != nil {
		return nil, fmt.Errorf("failed to query table: %w", err)
	}
	return &entity, nil
}

func (r *repo) ListByCodebaseAndUser(codebaseID string, userID users.ID) ([]*view.View, error) {
	var views []*view.View
	err := r.db.Select(&views, "SELECT id, user_id, codebase_id, workspace_id, name, last_used_at, created_at, mount_path, mount_hostname FROM views WHERE codebase_id=$1 AND user_id=$2", codebaseID, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to query table: %w", err)
	}
	return views, nil
}

func (r *repo) ListByCodebaseAndWorkspace(codebaseID, workspaceID string) ([]*view.View, error) {
	var views []*view.View
	err := r.db.Select(&views, `SELECT id, user_id, codebase_id, workspace_id, name, last_used_at, created_at, mount_path, mount_hostname
		FROM views WHERE codebase_id=$1 AND workspace_id=$2`, codebaseID, workspaceID)
	if err != nil {
		return nil, fmt.Errorf("failed to query table: %w", err)
	}
	return views, nil
}

func (r *repo) Update(e *view.View) error {
	_, err := r.db.NamedExec(`UPDATE views
		SET workspace_id = :workspace_id,
		    name = :name,
		    last_used_at = :last_used_at,
		    mount_path = :mount_path,
		    mount_hostname = :mount_hostname
		WHERE id = :id`, e)
	if err != nil {
		return fmt.Errorf("update view failed: %w", err)
	}
	return nil
}
