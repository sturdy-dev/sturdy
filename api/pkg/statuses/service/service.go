package service

import (
	"context"
	"fmt"

	"getsturdy.com/api/pkg/statuses"
	db_statuses "getsturdy.com/api/pkg/statuses/db"
	"getsturdy.com/api/pkg/view/events"

	"go.uber.org/zap"
)

type Service struct {
	logger       *zap.Logger
	repo         *db_statuses.Repository
	eventsSender events.EventSender
}

func New(
	logger *zap.Logger,
	repo *db_statuses.Repository,
	eventsSender events.EventSender,
) *Service {
	return &Service{
		logger:       logger,
		repo:         repo,
		eventsSender: eventsSender,
	}
}

var ErrInvalidStatus = fmt.Errorf("invalid status")

func (s *Service) Set(ctx context.Context, status *statuses.Status) error {
	if !statuses.ValidType[status.Type] {
		return ErrInvalidStatus
	}
	if err := s.repo.Create(ctx, status); err != nil {
		return fmt.Errorf("failed to create status: %w", err)
	}
	if err := s.eventsSender.Codebase(status.CodebaseID, events.StatusUpdated, status.ID); err != nil {
		s.logger.Error("failed to send status updated event", zap.Error(err))
	}
	return nil
}

func (s *Service) Get(ctx context.Context, id string) (*statuses.Status, error) {
	return s.repo.Get(ctx, id)
}

func (s *Service) List(ctx context.Context, codebaseID, commitID string) ([]*statuses.Status, error) {
	return s.repo.ListByCodebaseIDAndCommitID(ctx, codebaseID, commitID)
}
