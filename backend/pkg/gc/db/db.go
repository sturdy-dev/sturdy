package db

import (
	"context"
	"fmt"
	"mash/pkg/gc"
	"time"

	"github.com/jmoiron/sqlx"
)

type Repository interface {
	ListSince(ctx context.Context, codebaseID string, since time.Time) ([]*gc.CodebaseGarbageStatus, error)
	Create(context.Context, *gc.CodebaseGarbageStatus) error
}

type repo struct {
	db *sqlx.DB
}

func NewRepository(db *sqlx.DB) Repository {
	return &repo{db: db}
}

func (r *repo) ListSince(ctx context.Context, codebaseID string, since time.Time) ([]*gc.CodebaseGarbageStatus, error) {
	var res []*gc.CodebaseGarbageStatus
	if err := r.db.SelectContext(ctx, &res, `
		SELECT
			codebase_id,
			completed_at,
			duration_millis
		FROM
			codebases_garbage_collection_status
		WHERE
			codebase_id = $1
			AND completed_at > $2
	`, codebaseID, since); err != nil {
		return nil, fmt.Errorf("failed to select: %w", err)
	}
	return res, nil
}

func (r *repo) Create(ctx context.Context, status *gc.CodebaseGarbageStatus) error {
	if _, err := r.db.NamedExecContext(ctx, `
		INSERT INTO codebases_garbage_collection_status 
			(codebase_id, completed_at, duration_millis)
		VALUES
			(:codebase_id, :completed_at, :duration_millis)
	`, status); err != nil {
		return fmt.Errorf("failed to insert: %w", err)
	}
	return nil
}

func (r *repo) GetByCodebaseID(ctx context.Context, codebaseID string) (*gc.CodebaseGarbageStatus, error) {
	var res gc.CodebaseGarbageStatus
	if err := r.db.Get(&res, `
		SELECT 
			codebase_id,
			completed_at,
			duration_millis
		FROM 
			codebases_garbage_collection_status
		WHERE 
			codebase_id = $1
	`, codebaseID); err != nil {
		return nil, fmt.Errorf("failed to get: %w", err)
	}
	return &res, nil
}
