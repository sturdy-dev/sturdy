package service

import (
	"context"
	"fmt"

	"getsturdy.com/api/pkg/codebases"
	"getsturdy.com/api/pkg/events/v2"
	"getsturdy.com/api/pkg/statuses"
	db_statuses "getsturdy.com/api/pkg/statuses/db"

	"go.uber.org/zap"
)

type Service struct {
	logger          *zap.Logger
	repo            *db_statuses.Repository
	eventsPublisher *events.Publisher
}

func New(
	logger *zap.Logger,
	repo *db_statuses.Repository,
	eventsPublisher *events.Publisher,
) *Service {
	return &Service{
		logger:          logger,
		repo:            repo,
		eventsPublisher: eventsPublisher,
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
	if err := s.eventsPublisher.StatusUpdated(ctx, events.Codebase(status.CodebaseID), status); err != nil {
		s.logger.Error("failed to send status updated event", zap.Error(err))
	}
	return nil
}

func (s *Service) Get(ctx context.Context, id string) (*statuses.Status, error) {
	return s.repo.Get(ctx, id)
}

func (s *Service) List(ctx context.Context, codebaseID codebases.ID, commitID string) ([]*statuses.Status, error) {
	return s.repo.ListByCodebaseIDAndCommitID(ctx, codebaseID, commitID)
}

func (s *Service) ListByWorkspaceID(ctx context.Context, workspaceID string) ([]*statuses.Status, error) {
	return s.repo.ListByWorkspaceID(ctx, workspaceID)
}
