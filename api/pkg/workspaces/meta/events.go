package meta

import (
	"context"
	"time"

	"getsturdy.com/api/pkg/codebases"
	"getsturdy.com/api/pkg/events"
	"getsturdy.com/api/pkg/workspaces"
	"getsturdy.com/api/pkg/workspaces/db"

	"go.uber.org/zap"
)

type writerWithEvents struct {
	logger        *zap.Logger
	workspaceRepo db.Repository
	eventSender   events.EventSender
}

func NewWriterWithEvents(logger *zap.Logger, workspaceRepo db.Repository, eventSender events.EventSender) db.WorkspaceWriter {
	return &writerWithEvents{
		logger:        logger,
		eventSender:   eventSender,
		workspaceRepo: workspaceRepo,
	}
}

func (w *writerWithEvents) Create(workspace workspaces.Workspace) error {
	if err := w.workspaceRepo.Create(workspace); err != nil {
		return err
	}
	if err := w.eventSender.Codebase(workspace.CodebaseID, events.WorkspaceUpdated, workspace.ID); err != nil {
		w.logger.Error("failed to send event: %v", zap.Error(err))
		// do not fail
	}
	return nil
}

func (w *writerWithEvents) UpdateFields(ctx context.Context, workspaceID string, fields ...db.UpdateOption) error {
	if err := w.workspaceRepo.UpdateFields(ctx, workspaceID, fields...); err != nil {
		return err
	}
	ws, err := w.workspaceRepo.Get(workspaceID)
	if err != nil {
		return err
	}
	if err := w.eventSender.Codebase(ws.CodebaseID, events.WorkspaceUpdated, ws.ID); err != nil {
		w.logger.Error("failed to send event: %v", zap.Error(err))
		// do not fail
	}
	return nil
}

func (w *writerWithEvents) UnsetUpToDateWithTrunkForAllInCodebase(codebaseID codebases.ID) error {
	err := w.workspaceRepo.UnsetUpToDateWithTrunkForAllInCodebase(codebaseID)
	if err != nil {
		return err
	}
	workspaces, err := w.workspaceRepo.ListByCodebaseIDs([]codebases.ID{codebaseID}, false)
	if err != nil {
		return err
	}
	for _, ws := range workspaces {
		if err := w.eventSender.Codebase(ws.CodebaseID, events.WorkspaceUpdated, ws.ID); err != nil {
			w.logger.Error("failed to send event: %v", zap.Error(err))
			// do not fail
		}
	}
	return nil
}

// Updated sets UpdatedAt, and resets Behind and Ahead counters
func Updated(ctx context.Context, workspaceReader db.WorkspaceReader, workspaceWriter db.WorkspaceWriter, workspaceID string) error {
	now := time.Now()
	if err := workspaceWriter.UpdateFields(ctx, workspaceID,
		db.SetUpdatedAt(&now),
		// Is recalculated on next get/list
		db.SetUpToDateWithTrunk(nil),
		db.SetHeadChangeID(nil),
		db.SetHeadChangeComputed(false),
	); err != nil {
		return err
	}
	return nil
}
