package service

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"mash/pkg/view/events"
	"mash/pkg/workspace/activity"
	db_activity "mash/pkg/workspace/activity/db"

	"github.com/google/uuid"
)

type Service struct {
	workspaceActivityReadsRepo db_activity.ActivityReadsRepository
	eventsSender               events.EventSender
}

func New(workspaceActivityReadsRepo db_activity.ActivityReadsRepository,
	eventsSender events.EventSender) *Service {
	return &Service{
		workspaceActivityReadsRepo: workspaceActivityReadsRepo,
		eventsSender:               eventsSender,
	}
}

func (svc *Service) MarkAsRead(ctx context.Context, userID string, act *activity.WorkspaceActivity) error {
	lastRead, err := svc.workspaceActivityReadsRepo.GetByUserAndWorkspace(ctx, userID, act.WorkspaceID)
	switch {
	case err == nil:
	case errors.Is(err, sql.ErrNoRows):
		lastRead = &activity.WorkspaceActivityReads{
			ID:                uuid.NewString(),
			UserID:            userID,
			WorkspaceID:       act.WorkspaceID,
			LastReadCreatedAt: act.CreatedAt,
		}
		// create new
		if err := svc.workspaceActivityReadsRepo.Create(ctx, *lastRead); err != nil {
			return err
		}
	default:
		return fmt.Errorf("failed to get last read: %w", err)
	}

	// Update
	if lastRead.LastReadCreatedAt.Before(act.CreatedAt) {
		lastRead.LastReadCreatedAt = act.CreatedAt
		if err := svc.workspaceActivityReadsRepo.Update(ctx, lastRead); err != nil {
			return err
		}
	}

	// Send event (send to self that it has been read)
	svc.eventsSender.User(userID, events.WorkspaceUpdatedActivity, act.ID)

	return nil
}
