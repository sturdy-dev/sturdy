package service

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"getsturdy.com/api/pkg/activity"
	db_activity "getsturdy.com/api/pkg/activity/db"
	"getsturdy.com/api/pkg/changes"
	"getsturdy.com/api/pkg/events"
	"getsturdy.com/api/pkg/users"

	"github.com/google/uuid"
)

type Service struct {
	readsRepo    db_activity.ActivityReadsRepository
	repo         db_activity.ActivityRepository
	eventsSender events.EventSender
}

func New(
	readsRepo db_activity.ActivityReadsRepository,
	repo db_activity.ActivityRepository,
	eventsSender events.EventSender,
) *Service {
	return &Service{
		readsRepo:    readsRepo,
		repo:         repo,
		eventsSender: eventsSender,
	}
}

// SetChange updated activities for a given workspace with change_id where it's not set.
// This is used to "snapshot" change activities coming from a workspace and make them change specific.
func (svc *Service) SetChange(ctx context.Context, workspaceID string, changeID changes.ID) error {
	return svc.repo.SetChangeID(ctx, workspaceID, changeID)
}

func safeLimit(limit *int32) int32 {
	const maxLimit = 100
	const defaultLimit = 25
	if limit == nil {
		return defaultLimit
	}
	if *limit > maxLimit {
		return maxLimit
	}
	return *limit
}

func (svc *Service) ListByChangeID(ctx context.Context, changeID changes.ID, limit *int32) ([]*activity.Activity, error) {
	return svc.repo.ListByChangeID(ctx, changeID, safeLimit(limit))
}

func (svc *Service) ListByWorkspaceID(ctx context.Context, workspaceID string, limit *int32, after *time.Time) ([]*activity.Activity, error) {
	if after == nil {
		return svc.repo.ListByWorkspaceID(ctx, workspaceID, safeLimit(limit))
	}
	return svc.repo.ListByWorkspaceIDNewerThan(ctx, workspaceID, *after, safeLimit(limit))
}

func (svc *Service) MarkAsRead(ctx context.Context, userID users.ID, act *activity.Activity) error {
	if act.WorkspaceID == nil {
		return nil
	}

	lastRead, err := svc.readsRepo.GetByUserAndWorkspace(ctx, userID, *act.WorkspaceID)
	switch {
	case err == nil:
	case errors.Is(err, sql.ErrNoRows):
		lastRead = &activity.ActivityReads{
			ID:                uuid.NewString(),
			UserID:            userID,
			WorkspaceID:       *act.WorkspaceID,
			LastReadCreatedAt: act.CreatedAt,
		}
		// create new
		if err := svc.readsRepo.Create(ctx, *lastRead); err != nil {
			return err
		}
	default:
		return fmt.Errorf("failed to get last read: %w", err)
	}

	// Update
	if lastRead.LastReadCreatedAt.Before(act.CreatedAt) {
		lastRead.LastReadCreatedAt = act.CreatedAt
		if err := svc.readsRepo.Update(ctx, lastRead); err != nil {
			return err
		}
	}

	// Send event (send to self that it has been read)
	svc.eventsSender.User(userID, events.WorkspaceUpdatedActivity, act.ID)

	return nil
}
