package service

import (
	"context"
	"fmt"
	"sort"

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
	ss, err := s.repo.ListByWorkspaceID(ctx, workspaceID)
	if err != nil {
		return nil, err
	}
	if len(ss) == 0 {
		return nil, nil
	}

	// find latest status
	sort.Slice(ss, func(i, j int) bool {
		return ss[i].Timestamp.After(ss[j].Timestamp)
	})

	// return only statuses from the latest commit
	latestStatuses := make([]*statuses.Status, 0, len(ss))
	for _, s := range ss {
		if s.CommitSHA == ss[0].CommitSHA {
			latestStatuses = append(latestStatuses, s)
		}
	}
	return latestStatuses, nil
}

func (s *Service) NotifyAllInWorkspace(ctx context.Context, workspaceID string) error {
	statusList, err := s.ListByWorkspaceID(ctx, workspaceID)
	if err != nil {
		return fmt.Errorf("failed to get statuses: %w", err)
	}

	for _, status := range statusList {
		if err := s.eventsPublisher.StatusUpdated(ctx, events.Codebase(status.CodebaseID), status); err != nil {
			s.logger.Error("failed to send status updated event", zap.Error(err))
		}
	}

	return nil
}
