package service

import (
	"context"
	"fmt"
	"time"

	"getsturdy.com/api/pkg/events"
	"getsturdy.com/api/pkg/workspaces/watchers"
	db_watchers "getsturdy.com/api/pkg/workspaces/watchers/db"
)

type Service struct {
	repo         db_watchers.Repository
	eventsSender events.EventSender
}

func New(
	repo db_watchers.Repository,
	eventsSender events.EventSender,
) *Service {
	return &Service{
		repo:         repo,
		eventsSender: eventsSender,
	}
}

func (s *Service) Watch(ctx context.Context, userID, workspaceID string) (*watchers.Watcher, error) {
	watcher := &watchers.Watcher{
		UserID:      userID,
		WorkspaceID: workspaceID,
		Status:      watchers.StatusWatching,
		CreatedAt:   time.Now(),
	}
	if err := s.repo.Create(ctx, watcher); err != nil {
		return nil, fmt.Errorf("failed to create watcher: %w", err)
	}
	s.eventsSender.User(userID, events.WorkspaceWatchingStatusUpdated, workspaceID)
	return watcher, nil
}

func (s *Service) Unwatch(ctx context.Context, userID, workspaceID string) (*watchers.Watcher, error) {
	watcher := &watchers.Watcher{
		UserID:      userID,
		WorkspaceID: workspaceID,
		Status:      watchers.StatusIgnored,
		CreatedAt:   time.Now(),
	}
	if err := s.repo.Create(ctx, watcher); err != nil {
		return nil, fmt.Errorf("failed to create watcher: %w", err)
	}
	s.eventsSender.User(userID, events.WorkspaceWatchingStatusUpdated, workspaceID)
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
func (s *Service) Get(ctx context.Context, userID, workspaceID string) (*watchers.Watcher, error) {
	watcher, err := s.repo.GetByUserIDAndWorkspaceID(ctx, userID, workspaceID)
	if err != nil {
		return nil, fmt.Errorf("failed to get watcher: %w", err)
	}
	return watcher, nil
}
