package sender

import (
	"context"
	"fmt"
	"time"

	"getsturdy.com/api/pkg/activity"
	db_activity "getsturdy.com/api/pkg/activity/db"
	service_activity "getsturdy.com/api/pkg/activity/service"
	"getsturdy.com/api/pkg/auth"
	"getsturdy.com/api/pkg/codebases"
	db_codebases "getsturdy.com/api/pkg/codebases/db"
	"getsturdy.com/api/pkg/comments"
	"getsturdy.com/api/pkg/events"
	"getsturdy.com/api/pkg/users"

	"github.com/google/uuid"
)

type ActivitySender interface {
	// Codebase is deprecated. Create activity-type specific methods insteaad, like Comment
	Codebase(ctx context.Context, codebaseID codebases.ID, workspaceID string, userID users.ID, activityType activity.Type, referenceID string) error

	Comment(context.Context, *comments.Comment) error
}

type realActivitySender struct {
	codebaseUserRepo      db_codebases.CodebaseUserRepository
	workspaceActivityRepo db_activity.ActivityRepository
	activityService       *service_activity.Service
	eventsSender          events.EventSender
}

func NewActivitySender(
	codebaseUserRepo db_codebases.CodebaseUserRepository,
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

func (s *realActivitySender) Comment(ctx context.Context, comment *comments.Comment) error {
	userID, err := auth.UserID(ctx)
	if err != nil {
		return fmt.Errorf("failed to get user id from context: %w", err)
	}

	act := activity.Activity{
		ID:           uuid.NewString(),
		UserID:       userID,
		CreatedAt:    time.Now(),
		ActivityType: activity.TypeComment,
		Reference:    comment.ID.String(),
		WorkspaceID:  comment.WorkspaceID,
		ChangeID:     comment.ChangeID,
	}

	if err := s.workspaceActivityRepo.Create(ctx, act); err != nil {
		return fmt.Errorf("failed to create activity: %w", err)
	}

	// Send to all members of this codebase
	codebaseUsers, err := s.codebaseUserRepo.GetByCodebase(comment.CodebaseID)
	if err != nil {
		return fmt.Errorf("failed to get codebase users: %w", err)
	}

	for _, codebaseUser := range codebaseUsers {
		s.eventsSender.User(codebaseUser.UserID, events.WorkspaceUpdatedActivity, act.ID)
	}

	return nil
}

func (s *realActivitySender) Codebase(ctx context.Context, codebaseID codebases.ID, workspaceID string, userID users.ID, activityType activity.Type, referenceID string) error {
	activityID := uuid.NewString()

	act := activity.Activity{
		ID:           activityID,
		UserID:       userID,
		WorkspaceID:  &workspaceID,
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

func (noopActivitySender) Codebase(ctx context.Context, codebaseID codebases.ID, workspaceID string, userID users.ID, activityType activity.Type, referenceID string) error {
	return nil
}

func (noopActivitySender) Comment(context.Context, *comments.Comment) error {
	return nil
}

func NewNoopNotificationSender() ActivitySender {
	return noopActivitySender{}
}
