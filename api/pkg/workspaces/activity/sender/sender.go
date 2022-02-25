package sender

import (
	"context"
	"time"

	db_codebase "getsturdy.com/api/pkg/codebase/db"
	"getsturdy.com/api/pkg/events"
	"getsturdy.com/api/pkg/users"
	"getsturdy.com/api/pkg/workspaces/activity"
	db_activity "getsturdy.com/api/pkg/workspaces/activity/db"
	service_activity "getsturdy.com/api/pkg/workspaces/activity/service"

	"github.com/google/uuid"
)

type ActivitySender interface {
	Codebase(ctx context.Context, codebaseID, workspaceID string, userID users.ID, activityType activity.WorkspaceActivityType, referenceID string) error
}

type realActivitySender struct {
	codebaseUserRepo      db_codebase.CodebaseUserRepository
	workspaceActivityRepo db_activity.ActivityRepository
	activityService       *service_activity.Service
	eventsSender          events.EventSender
}

func NewActivitySender(
	codebaseUserRepo db_codebase.CodebaseUserRepository,
	workspaceActivityRepo db_activity.ActivityRepository,
	activityService *service_activity.Service,
	eventsSender events.EventSender,
) ActivitySender {
	return &realActivitySender{
		codebaseUserRepo:      codebaseUserRepo,
		workspaceActivityRepo: workspaceActivityRepo,
		activityService:       activityService,
		eventsSender:          eventsSender,
	}
}

func (s *realActivitySender) Codebase(ctx context.Context, codebaseID, workspaceID string, userID users.ID, activityType activity.WorkspaceActivityType, referenceID string) error {
	activityID := uuid.NewString()

	act := activity.WorkspaceActivity{
		ID:           activityID,
		UserID:       userID,
		WorkspaceID:  workspaceID,
		CreatedAt:    time.Now(),
		ActivityType: activityType,
		Reference:    referenceID,
	}

	if err := s.workspaceActivityRepo.Create(ctx, act); err != nil {
		return err
	}

	// Before sending, mark as read for the sending user
	if err := s.activityService.MarkAsRead(ctx, userID, &act); err != nil {
		return err
	}

	// Send to all members of this codebase
	codebaseUsers, err := s.codebaseUserRepo.GetByCodebase(codebaseID)
	if err != nil {
		return err
	}
	for _, codebaseUser := range codebaseUsers {
		s.eventsSender.User(codebaseUser.UserID, events.WorkspaceUpdatedActivity, activityID)
	}

	return nil
}

type noopActivitySender struct{}

func (noopActivitySender) Codebase(ctx context.Context, codebaseID, workspaceID string, userID users.ID, activityType activity.WorkspaceActivityType, referenceID string) error {
	return nil
}

func NewNoopNotificationSender() ActivitySender {
	return noopActivitySender{}
}
