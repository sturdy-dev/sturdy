package db

import (
	"context"
	"fmt"
	"getsturdy.com/api/pkg/review"
	"time"

	"github.com/jmoiron/sqlx"
)

type database struct {
	db *sqlx.DB
}

func NewReviewRepository(db *sqlx.DB) ReviewRepository {
	return &database{db: db}
}

func (r *database) Create(ctx context.Context, rev review.Review) error {
	_, err := r.db.NamedExecContext(ctx, `INSERT INTO workspace_reviews (id, codebase_id, workspace_id, user_id, grade, created_at, is_replaced, requested_by)
		VALUES(:id, :codebase_id, :workspace_id, :user_id, :grade, :created_at, :is_replaced, :requested_by)`, rev)
	if err != nil {
		return fmt.Errorf("failed to insert review: %w", err)
	}
	return nil
}

func (r *database) Update(ctx context.Context, rev *review.Review) error {
	_, err := r.db.NamedExecContext(ctx, `UPDATE workspace_reviews
		SET grade = :grade,
		    dismissed_at = :dismissed_at,
		    is_replaced = :is_replaced,
		    requested_by = :requested_by
		WHERE id = :id`, rev)
	if err != nil {
		return fmt.Errorf("failed to insert review: %w", err)
	}
	return nil
}

func (r *database) Get(ctx context.Context, id string) (*review.Review, error) {
	var res review.Review
	err := r.db.GetContext(ctx, &res, `SELECT id, codebase_id, workspace_id, user_id, grade, created_at, dismissed_at, is_replaced, requested_by
		FROM workspace_reviews
		WHERE id = $1`, id)
	if err != nil {
		return nil, fmt.Errorf("failed to list by workspace: %w", err)
	}
	return &res, nil
}

func (r *database) GetLatestByUserAndWorkspace(ctx context.Context, userID, workspaceID string) (*review.Review, error) {
	var res review.Review
	err := r.db.GetContext(ctx, &res, `SELECT id, codebase_id, workspace_id, user_id, grade, created_at, dismissed_at, is_replaced, requested_by
		FROM workspace_reviews
		WHERE workspace_id = $1
	      AND user_id = $2
	      AND is_replaced IS FALSE`, workspaceID, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to list by workspace: %w", err)
	}
	return &res, nil
}

func (r *database) ListLatestByWorkspace(ctx context.Context, workspaceID string) ([]*review.Review, error) {
	var res []*review.Review
	err := r.db.SelectContext(ctx, &res, `SELECT id, codebase_id, workspace_id, user_id, grade, created_at, dismissed_at, is_replaced, requested_by
		FROM workspace_reviews
		WHERE workspace_id = $1
		AND dismissed_at IS NULL
		AND is_replaced IS FALSE`, workspaceID)
	if err != nil {
		return nil, fmt.Errorf("failed to list by workspace: %w", err)
	}
	return res, nil
}

func (r *database) DismissAllInWorkspace(ctx context.Context, workspaceID string) error {
	ts := time.Now()
	_, err := r.db.ExecContext(ctx, `UPDATE workspace_reviews
		SET dismissed_at = $1
		WHERE workspace_id = $2 
			AND dismissed_at IS NULL
			AND is_replaced IS FALSE`, ts, workspaceID)
	if err != nil {
		return fmt.Errorf("failed to list by workspace: %w", err)
	}
	return nil
}
