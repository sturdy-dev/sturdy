package db

import (
	"context"
	"fmt"

	"mash/pkg/suggestions"

	"github.com/jmoiron/sqlx"
)

var _ Repository = &database{}

type database struct {
	db *sqlx.DB
}

func New(d *sqlx.DB) Repository {
	return &database{
		db: d,
	}
}

func (d *database) Create(ctx context.Context, suggestion *suggestions.Suggestion) error {
	if _, err := d.db.NamedExecContext(ctx, `
		INSERT INTO suggestions_v2 (
			id,
			codebase_id,
			workspace_id,
			for_workspace_id,
			for_snapshot_id,
			created_at,
			applied_hunks,
			dismissed_hunks,
			user_id,
		    dismissed_at,
		    notified_at
		) VALUES (
			:id,
			:codebase_id,
			:workspace_id,
			:for_workspace_id,
			:for_snapshot_id,
			:created_at,
			:applied_hunks,
			:dismissed_hunks,
			:user_id,
		    :dismissed_at,
		    :notified_at
		)`, suggestion); err != nil {
		return fmt.Errorf("failed to insert: %w", err)
	}
	return nil
}

func (d *database) Update(ctx context.Context, suggestion *suggestions.Suggestion) error {
	if _, err := d.db.NamedExecContext(ctx, `
		UPDATE suggestions_v2
		SET 
			applied_hunks = :applied_hunks,
			dismissed_hunks = :dismissed_hunks,
		    dismissed_at = :dismissed_at,
		    notified_at = :notified_at
		WHERE 
			id = :id`, suggestion); err != nil {
		return fmt.Errorf("failed to update: %w", err)
	}
	return nil
}

func (d *database) GetByID(ctx context.Context, id suggestions.ID) (*suggestions.Suggestion, error) {
	suggestion := &suggestions.Suggestion{}
	if err := d.db.GetContext(ctx, suggestion, `
		SELECT
			id,
			codebase_id,
			workspace_id,
			for_workspace_id,
			for_snapshot_id,
			created_at,
			applied_hunks,
			dismissed_hunks,
			user_id,
			dismissed_at,
			notified_at
		FROM suggestions_v2
		WHERE id = $1`, id); err != nil {
		return nil, fmt.Errorf("failed to get suggestion: %w", err)
	}
	return suggestion, nil
}

func (d *database) GetByWorkspaceID(ctx context.Context, id string) (*suggestions.Suggestion, error) {
	suggestion := &suggestions.Suggestion{}
	if err := d.db.GetContext(ctx, suggestion, `
		SELECT
			id,
			codebase_id,
			workspace_id,
			for_workspace_id,
			for_snapshot_id,
			created_at,
			applied_hunks,
			dismissed_hunks,
			user_id,
			dismissed_at,
			notified_at
		FROM suggestions_v2
		WHERE workspace_id = $1`, id); err != nil {
		return nil, fmt.Errorf("failed to get suggestion: %w", err)
	}
	return suggestion, nil
}

func (d *database) ListForWorkspaceID(ctx context.Context, id string) ([]*suggestions.Suggestion, error) {
	suggestions := []*suggestions.Suggestion{}
	if err := d.db.SelectContext(ctx, &suggestions, `
		SELECT
			id,
			codebase_id,
			workspace_id,
			for_workspace_id,
			for_snapshot_id,
			created_at,
			applied_hunks,
			dismissed_hunks,
			user_id,
			dismissed_at,
			notified_at
		FROM suggestions_v2
		WHERE for_workspace_id = $1`, id); err != nil {
		return nil, fmt.Errorf("failed to get suggestions: %w", err)
	}
	return suggestions, nil
}

func (d *database) ListBySnapshotID(ctx context.Context, snapshotID string) ([]*suggestions.Suggestion, error) {
	suggestions := []*suggestions.Suggestion{}
	if err := d.db.SelectContext(ctx, &suggestions, `
		SELECT
			id,
			codebase_id,
			workspace_id,
			for_workspace_id,
			for_snapshot_id,
			created_at,
			applied_hunks,
			dismissed_hunks,
			user_id,
			dismissed_at,
			notified_at
		FROM suggestions_v2
		WHERE for_snapshot_id = $1`, snapshotID); err != nil {
		return nil, fmt.Errorf("failed to get suggestions: %w", err)
	}
	return suggestions, nil
}
