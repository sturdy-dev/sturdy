package service

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/pkg/errors"

	events "getsturdy.com/api/pkg/events/v2"
	"getsturdy.com/api/pkg/snapshots"
	db_snapshots "getsturdy.com/api/pkg/snapshots/db"
	service_snapshots "getsturdy.com/api/pkg/snapshots/service"
	"getsturdy.com/api/pkg/view"
	"getsturdy.com/api/pkg/view/db"
	vcs_view "getsturdy.com/api/pkg/view/vcs"
	"getsturdy.com/api/pkg/workspaces"
	db_workspaces "getsturdy.com/api/pkg/workspaces/db"
	"getsturdy.com/api/vcs"
	"getsturdy.com/api/vcs/executor"

	"go.uber.org/zap"
)

type Service struct {
	logger           *zap.Logger
	viewRepo         db.Repository
	workspaceReader  db_workspaces.WorkspaceReader
	gitSnapshotter   *service_snapshots.Service
	snapshotRepo     db_snapshots.Repository
	workspaceWriter  db_workspaces.WorkspaceWriter
	executorProvider executor.Provider
	eventSender      *events.Publisher
}

func New(
	logger *zap.Logger,
	viewRepo db.Repository,
	workspaceReader db_workspaces.WorkspaceReader,
	gitSnapshotter *service_snapshots.Service,
	snapshotRepo db_snapshots.Repository,
	workspaceWriter db_workspaces.WorkspaceWriter,
	executorProvider executor.Provider,
	eventSender *events.Publisher,
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

var ErrRebasing = errors.New("is rebasing")

func (s *Service) OpenWorkspace(ctx context.Context, view *view.View, ws *workspaces.Workspace) error {
	if ws.ArchivedAt != nil {
		return fmt.Errorf("the workspace is archived")
	}
	if view.CodebaseID != ws.CodebaseID {
		return fmt.Errorf("the view and workspace does not belong to the same codebase")
	}
	if view.UserID != ws.UserID {
		return fmt.Errorf("the view and workspace does not belong to the same user")
	}

	// If the view that we're opening this workspace on has another workspace currently open, snapshot and save the changes
	currentWorkspaceOnView, err := s.workspaceReader.GetByViewID(view.ID, true)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return fmt.Errorf("could not find previous workspace on view: %w", err)
	} else if err == nil {
		_, err := s.gitSnapshotter.Snapshot(currentWorkspaceOnView.CodebaseID, currentWorkspaceOnView.ID,
			snapshots.ActionPreCheckoutOtherWorkspace, service_snapshots.WithOnView(*currentWorkspaceOnView.ViewID), service_snapshots.WithMarkAsLatestInWorkspace())
		if errors.Is(err, service_snapshots.ErrCantSnapshotRebasing) {
			return ErrRebasing
		}
		if err != nil {
			return fmt.Errorf("failed to snapshot: %w", err)
		}

		currentWorkspaceOnView.ViewID = nil
		if err := s.workspaceWriter.UpdateFields(ctx, currentWorkspaceOnView.ID, db_workspaces.SetViewID(nil)); err != nil {
			return fmt.Errorf("failed to finalize previous workspace on view: %w", err)
		}
	}

	// If the workspace that we're opening on the view, currently is opened somewhere else, snapshot it and save the changes
	// The snapshot will also be used to "move" the contents to the new view
	if ws.ViewID != nil {
		_, err := s.gitSnapshotter.Snapshot(ws.CodebaseID, ws.ID, snapshots.ActionPreCheckoutOtherView,
			service_snapshots.WithOnView(*ws.ViewID), service_snapshots.WithMarkAsLatestInWorkspace())
		if err != nil {
			return fmt.Errorf("failed to snapshot: %w", err)
		}

		// TODO: unset view.workspace_id?
	}

	if err := s.executorProvider.New().
		Write(vcs_view.CheckoutBranch(ws.ID)).
		Write(func(repo vcs.RepoWriter) error {
			// Restore snapshot
			if ws.LatestSnapshotID != nil {
				snapshot, err := s.snapshotRepo.Get(*ws.LatestSnapshotID)
				if err != nil {
					return fmt.Errorf("failed to get snapshot: %w", err)
				}
				if err := s.gitSnapshotter.Restore(snapshot, repo); err != nil {
					return fmt.Errorf("failed to restore snapshot: %w", err)
				}
			} else {
				// vcs2.Restore does this as well, make sure to smudge files in this scenario as well
				if err := repo.LargeFilesPull(); err != nil {
					s.logger.Warn("failed to pull large files", zap.Error(err))
				}
				if err := repo.CleanStaged(); err != nil {
					return fmt.Errorf("failed to clean index: %w", err)
				}
			}
			return nil
		}).ExecView(ws.CodebaseID, view.ID, "OpenWorkspaceOnView"); err != nil {
		return fmt.Errorf("failed to open workspace on view: %w", err)
	}

	// Update workspace object
	ws.ViewID = &view.ID
	if err := s.workspaceWriter.UpdateFields(ctx, ws.ID, db_workspaces.SetViewID(&view.ID)); err != nil {
		return fmt.Errorf("failed to update workspace: %w", err)
	}

	// Update view object
	view.WorkspaceID = ws.ID
	if err := s.viewRepo.Update(view); err != nil {
		return fmt.Errorf("failed to update view obj: %w", err)
	}

	if err := s.eventSender.ViewUpdated(ctx, events.Workspace(ws.ID), view); err != nil {
		s.logger.Error("failed to send view updated event: %w", zap.Error(err))
		// do not fail
	}

	return nil
}

func (s *Service) GetByID(_ context.Context, id string) (*view.View, error) {
	return s.viewRepo.Get(id)
}
