package service

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"getsturdy.com/api/pkg/events"
	"getsturdy.com/api/pkg/presence"
	db_presence "getsturdy.com/api/pkg/presence/db"
	"getsturdy.com/api/pkg/users"

	"github.com/google/uuid"
)

type Service interface {
	Record(ctx context.Context, userID users.ID, workspaceID string, state presence.State) (*presence.Presence, error)
	ListByWorkspace(ctx context.Context, workspaceID string) ([]*presence.Presence, error)
}

type service struct {
	presenceRepo db_presence.PresenceRepository
	eventSender  events.EventSender
}

func New(presenceRepo db_presence.PresenceRepository, eventSender events.EventSender) Service {
	return &service{
		presenceRepo: presenceRepo,
		eventSender:  eventSender,
	}
}

func (p *service) Record(ctx context.Context, userID users.ID, workspaceID string, state presence.State) (*presence.Presence, error) {
	pre, err := p.presenceRepo.GetByUserAndWorkspace(ctx, userID, workspaceID)
	if errors.Is(err, sql.ErrNoRows) {
		newPresence := &presence.Presence{
			ID:           uuid.NewString(),
			UserID:       userID,
			WorkspaceID:  workspaceID,
			LastActiveAt: time.Now(),
			State:        state,
		}

		// There might be a data race here, hence - upsert.
		if err := p.presenceRepo.Upsert(ctx, newPresence); err != nil {
			return nil, fmt.Errorf("failed to create new presence: %w", err)
		}

		// Send update
		if err := p.eventSender.Workspace(workspaceID, events.WorkspaceUpdatedPresence, workspaceID); err != nil {
			return nil, fmt.Errorf("failed to send presence event: %w", err)
		}

		return newPresence, nil
	} else if err != nil {
		return nil, fmt.Errorf("failed to record presence: %w", err)
	}

	// Don't downgrade to a "lesser" priority, unless that activity is more than 10 minutes ago
	if presence.StatePriority[pre.State] > presence.StatePriority[state] &&
		time.Now().Add(-1*time.Minute*10).Before(pre.LastActiveAt) {
		// Don't update
		return pre, nil
	}

	// Update existing
	pre.LastActiveAt = time.Now()
	pre.State = state

	if err := p.presenceRepo.Update(ctx, pre); err != nil {
		return nil, fmt.Errorf("failed to update presence: %w", err)
	}

	// Send update
	if err := p.eventSender.Workspace(workspaceID, events.WorkspaceUpdatedPresence, workspaceID); err != nil {
		return nil, fmt.Errorf("failed to send presence event: %w", err)
	}

	return pre, nil
}

func (p *service) ListByWorkspace(ctx context.Context, workspaceID string) ([]*presence.Presence, error) {
	return p.presenceRepo.ListByWorkspace(ctx, workspaceID)
}
