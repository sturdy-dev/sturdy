package service

import (
	"context"
	"fmt"

	"getsturdy.com/api/pkg/events"
	db_snapshots "getsturdy.com/api/pkg/snapshots/db"
	"getsturdy.com/api/pkg/snapshots/snapshotter"
	"getsturdy.com/api/pkg/view"
	"getsturdy.com/api/pkg/view/db"
	"getsturdy.com/api/pkg/view/open"
	"getsturdy.com/api/pkg/workspaces"
	db_workspaces "getsturdy.com/api/pkg/workspaces/db"
	"getsturdy.com/api/vcs/executor"

	"go.uber.org/zap"
)

type Service struct {
	logger           *zap.Logger
	viewRepo         db.Repository
	workspaceReader  db_workspaces.WorkspaceReader
	gitSnapshotter   snapshotter.Snapshotter
	snapshotRepo     db_snapshots.Repository
	workspaceWriter  db_workspaces.WorkspaceWriter
	executorProvider executor.Provider
	eventSender      events.EventSender
}

func New(
	logger *zap.Logger,
	viewRepo db.Repository,
	workspaceReader db_workspaces.WorkspaceReader,
	gitSnapshotter snapshotter.Snapshotter,
	snapshotRepo db_snapshots.Repository,
	workspaceWriter db_workspaces.WorkspaceWriter,
	executorProvider executor.Provider,
	eventSender events.EventSender,
) *Service {
	return &Service{
		logger:           logger.Named("views_service"),
		viewRepo:         viewRepo,
		workspaceReader:  workspaceReader,
		gitSnapshotter:   gitSnapshotter,
		snapshotRepo:     snapshotRepo,
		workspaceWriter:  workspaceWriter,
		executorProvider: executorProvider,
		eventSender:      eventSender,
	}
}

func (s *Service) OpenWorkspace(ctx context.Context, view *view.View, ws *workspaces.Workspace) error {
	if err := open.OpenWorkspaceOnView(
		ctx,
		s.logger,
		view,
		ws,
		s.viewRepo,
		s.workspaceReader,
		s.gitSnapshotter,
		s.snapshotRepo,
		s.workspaceWriter,
		s.executorProvider,
		s.eventSender,
	); err != nil {
		return fmt.Errorf("failed to open workspace: %w", err)
	}
	return nil
}

func (s *Service) GetByID(_ context.Context, id string) (*view.View, error) {
	return s.viewRepo.Get(id)
}
