package sender

import (
	"context"
	"errors"
	"fmt"
	"time"

	db_codebase "getsturdy.com/api/pkg/codebase/db"
	"getsturdy.com/api/pkg/emails/transactional"
	"getsturdy.com/api/pkg/events"
	"getsturdy.com/api/pkg/notification"
	db_notification "getsturdy.com/api/pkg/notification/db"
	"getsturdy.com/api/pkg/users"
	db_user "getsturdy.com/api/pkg/users/db"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

type NotificationSender interface {
	Codebase(ctx context.Context, codebaseID string, notificationType notification.NotificationType, referenceID string, senderUserID users.ID) error
	User(ctx context.Context, userID users.ID, codebaseID string, notificationType notification.NotificationType, referenceID string) error
}

type realNotificationSender struct {
	logger *zap.Logger

	codebaseUserRepo db_codebase.CodebaseUserRepository
	notificationRepo db_notification.Repository
	userRepo         db_user.Repository

	eventsSender events.EventSender
	emailSender  transactional.EmailSender
}

func NewNotificationSender(
	logger *zap.Logger,

	codebaseUserRepo db_codebase.CodebaseUserRepository,
	notificationRepo db_notification.Repository,
	userRepo db_user.Repository,

	eventsSender events.EventSender,
	emailSender transactional.EmailSender,
) NotificationSender {
	return &realNotificationSender{
		logger: logger,

		codebaseUserRepo: codebaseUserRepo,
		notificationRepo: notificationRepo,
		userRepo:         userRepo,

		eventsSender: eventsSender,
		emailSender:  emailSender,
	}
}

func (s *realNotificationSender) Codebase(ctx context.Context, codebaseID string, notificationType notification.NotificationType, referenceID string, senderUserID users.ID) error {
	// Send to all members of this codebase
	codebaseUsers, err := s.codebaseUserRepo.GetByCodebase(codebaseID)
	if err != nil {
		return err
	}
	for _, codebaseUser := range codebaseUsers {
		// Don't send to yourself
		if codebaseUser.UserID == senderUserID {
			continue
		}

		notif := notification.Notification{
			ID:               uuid.NewString(),
			UserID:           codebaseUser.UserID,
			CodebaseID:       codebaseID,
			CreatedAt:        time.Now(),
			NotificationType: notificationType,
			ReferenceID:      referenceID,
		}
		if err := s.notificationRepo.Create(notif); err != nil {
			return err
		}

		if err := s.dispatch(ctx, &notif); err != nil {
			return fmt.Errorf("failed to dispatch notification: %w", err)
		}

	}
	return nil
}

func (s *realNotificationSender) User(ctx context.Context, userID users.ID, codebaseID string, notificationType notification.NotificationType, referenceID string) error {
	notif := notification.Notification{
		ID:               uuid.NewString(),
		UserID:           userID,
		CodebaseID:       codebaseID,
		CreatedAt:        time.Now(),
		NotificationType: notificationType,
		ReferenceID:      referenceID,
	}

	if err := s.notificationRepo.Create(notif); err != nil {
		return fmt.Errorf("failed to save notification to the db: %w", err)
	}

	if err := s.dispatch(ctx, &notif); err != nil {
		return fmt.Errorf("failed to dispatch notification: %w", err)
	}

	return nil
}

func (s *realNotificationSender) dispatch(ctx context.Context, notif *notification.Notification) error {
	s.eventsSender.User(notif.UserID, events.NotificationEvent, notif.ID)

	user, err := s.userRepo.Get(notif.UserID)
	if err != nil {
		return fmt.Errorf("failed to find user: %w", err)
	}

	if err := s.emailSender.SendNotification(ctx, user, notif); errors.Is(err, transactional.ErrNotSupported) {
		s.logger.Warn("email notification not supported", zap.String("type", string(notif.NotificationType)))
	} else if err != nil {
		return fmt.Errorf("failed to notify via email: %w", err)
	}
	return nil
}

type noopNotificationSender struct{}

func (noopNotificationSender) Codebase(_ context.Context, codebaseID string, notificationType notification.NotificationType, referenceID string, senderUserID users.ID) error {
	return nil
}

func (noopNotificationSender) User(_ context.Context, userID users.ID, codebaseID string, notificationType notification.NotificationType, referenceID string) error {
	return nil
}

func NewNoopNotificationSender() NotificationSender {
	return noopNotificationSender{}
}
