package service

import (
	"context"
	"fmt"
	"time"

	"getsturdy.com/api/pkg/events/v2"
	"getsturdy.com/api/pkg/users"
	"getsturdy.com/api/pkg/workspaces/watchers"
	db_watchers "getsturdy.com/api/pkg/workspaces/watchers/db"
	"go.uber.org/zap"
)

type Service struct {
	repo         db_watchers.Repository
	logger       *zap.Logger
	eventsSender *events.Publisher
}

func New(
	repo db_watchers.Repository,
	logger *zap.Logger,
	eventsSender *events.Publisher,
) *Service {
	return &Service{
		repo:         repo,
		logger:       logger.Named("watchers_service"),
		eventsSender: eventsSender,
	}
}

func (s *Service) Watch(ctx context.Context, userID users.ID, workspaceID string) (*watchers.Watcher, error) {
	watcher := &watchers.Watcher{
		UserID:      userID,
		WorkspaceID: workspaceID,
		Status:      watchers.StatusWatching,
		CreatedAt:   time.Now(),
	}
	if err := s.repo.Create(ctx, watcher); err != nil {
		return nil, fmt.Errorf("failed to create watcher: %w", err)
	}
	if err := s.eventsSender.WorkspaceWatchingStatusUpdated(ctx, events.User(userID), watcher); err != nil {
		s.logger.Error("failed to send workspace watching status updated event", zap.Error(err))
	}
	return watcher, nil
}

func (s *Service) Unwatch(ctx context.Context, userID users.ID, workspaceID string) (*watchers.Watcher, error) {
	watcher := &watchers.Watcher{
		UserID:      userID,
		WorkspaceID: workspaceID,
		Status:      watchers.StatusIgnored,
		CreatedAt:   time.Now(),
	}
	if err := s.repo.Create(ctx, watcher); err != nil {
		return nil, fmt.Errorf("failed to create watcher: %w", err)
	}
	if err := s.eventsSender.WorkspaceWatchingStatusUpdated(ctx, events.User(userID), watcher); err != nil {
		s.logger.Error("failed to send workspace watching status updated event", zap.Error(err))
	}

	return watcher, nil
}

// List returns a list of watchers for a given workspace.
func (s *Service) ListWatchers(ctx context.Context, workspaceID string) ([]*watchers.Watcher, error) {
	watchers, err := s.repo.ListWatchingByWorkspaceID(ctx, workspaceID)
	if err != nil {
		return nil, fmt.Errorf("failed to list watchers: %w", err)
	}
	return watchers, nil
}

// Get returns a watcher for a given workspace and user.
func (s *Service) Get(ctx context.Context, userID users.ID, workspaceID string) (*watchers.Watcher, error) {
	watcher, err := s.repo.GetByUserIDAndWorkspaceID(ctx, userID, workspaceID)
	if err != nil {
		return nil, fmt.Errorf("failed to get watcher: %w", err)
	}
	return watcher, nil
}
