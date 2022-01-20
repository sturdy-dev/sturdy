package db

import (
	"context"
	"fmt"
	"getsturdy.com/api/pkg/presence"

	"github.com/jmoiron/sqlx"
)

type PresenceRepository interface {
	GetByUserAndWorkspace(ctx context.Context, userID, workspaceID string) (*presence.Presence, error)
	ListByWorkspace(ctx context.Context, workspaceID string) ([]*presence.Presence, error)
	Create(ctx context.Context, p presence.Presence) error
	Update(ctx context.Context, p *presence.Presence) error
}

type repo struct {
	db *sqlx.DB
}

func NewRepo(db *sqlx.DB) PresenceRepository {
	return &repo{db: db}
}

func (r *repo) GetByUserAndWorkspace(ctx context.Context, userID, workspaceID string) (*presence.Presence, error) {
	var p presence.Presence
	if err := r.db.GetContext(ctx, &p, `SELECT id, user_id, workspace_id, last_active_at, state FROM presence WHERE user_id = $1 AND workspace_id = $2`, userID, workspaceID); err != nil {
		return nil, fmt.Errorf("could not GetByUserAndWorkspace: %w", err)
	}
	return &p, nil
}

func (r *repo) ListByWorkspace(ctx context.Context, workspaceID string) ([]*presence.Presence, error) {
	var res []*presence.Presence
	if err := r.db.SelectContext(ctx, &res, `SELECT id, user_id, workspace_id, last_active_at, state FROM presence WHERE workspace_id = $1`, workspaceID); err != nil {
		return nil, fmt.Errorf("could not ListByWorkspace: %w", err)
	}
	return res, nil
}

func (r *repo) Create(ctx context.Context, p presence.Presence) error {
	if _, err := r.db.NamedExecContext(ctx, `INSERT INTO presence (id, user_id, workspace_id, last_active_at, state) VALUES(:id, :user_id, :workspace_id, :last_active_at, :state)`, p); err != nil {
		return fmt.Errorf("failed to Create: %w", err)
	}
	return nil
}

func (r *repo) Update(ctx context.Context, p *presence.Presence) error {
	if _, err := r.db.NamedExecContext(ctx, `UPDATE presence SET last_active_at = :last_active_at, state = :state WHERE id = :id`, p); err != nil {
		return fmt.Errorf("failed to update: %w", err)
	}
	return nil
}
