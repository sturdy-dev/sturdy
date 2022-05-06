package service

import (
	"context"
	"fmt"

	service_snapshots "getsturdy.com/api/pkg/snapshots/service"
	"getsturdy.com/api/pkg/statuses"
	service_statuses "getsturdy.com/api/pkg/statuses/service"
	"getsturdy.com/api/pkg/workspaces"
)

type Service struct {
	statusesService  *service_statuses.Service
	snapshotsService *service_snapshots.Service
}

func New(
	statusesService *service_statuses.Service,
	snapshotsService *service_snapshots.Service,
) *Service {
	return &Service{
		statusesService:  statusesService,
		snapshotsService: snapshotsService,
	}
}

func (s *Service) HealthyStatus(ctx context.Context, ws *workspaces.Workspace) (bool, error) {
	statusList, err := s.statusesService.ListByWorkspaceID(ctx, ws.ID)
	if err != nil {
		return false, err
	}

	// no status => is unhealthy
	if len(statusList) == 0 {
		return false, nil
	}

	for _, status := range statusList {
		if status.Type != statuses.TypeHealthy {
			return false, nil
		}

		isStale, err := s.StatusIsStaleForWorkspace(ctx, ws, status)
		if err != nil {
			return false, err
		}

		if isStale {
			return false, nil
		}
	}

	// have statuses, and all statuses are healthy
	return true, nil
}

func (s *Service) StatusIsStaleForWorkspace(ctx context.Context, ws *workspaces.Workspace, status *statuses.Status) (bool, error) {
	snapshot, err := s.snapshotsService.GetByCommitSHA(ctx, status.CommitSHA)
	if err != nil {
		return false, fmt.Errorf("failed to get snapshot: %w", err)
	}
	if ws.LatestSnapshotID == nil {
		return false, nil
	}
	if snapshot.ID != *ws.LatestSnapshotID {
		return true, nil
	}
	return false, nil
}
