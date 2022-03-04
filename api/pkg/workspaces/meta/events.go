package meta

import (
	"context"
	"time"

	"getsturdy.com/api/pkg/changes"
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

func (w *writerWithEvents) Update(ctx context.Context, workspace *workspaces.Workspace) error {
	if err := w.workspaceRepo.Update(ctx, workspace); err != nil {
		return err
	}
	if err := w.eventSender.Codebase(workspace.CodebaseID, events.WorkspaceUpdated, workspace.ID); err != nil {
		w.logger.Error("failed to send event: %v", zap.Error(err))
		// do not fail
	}
	return nil
}

func (w *writerWithEvents) SetUpToDateWithTrunk(ctx context.Context, workspaceID string, upToDateWithTrunk bool) error {
	if err := w.workspaceRepo.SetUpToDateWithTrunk(ctx, workspaceID, upToDateWithTrunk); err != nil {
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

func (w *writerWithEvents) SetHeadChange(ctx context.Context, workspaceID string, changeID changes.ID) error {
	if err := w.workspaceRepo.SetHeadChange(ctx, workspaceID, changeID); err != nil {
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

func (w *writerWithEvents) UnsetUpToDateWithTrunkForAllInCodebase(codebaseID string) error {
	err := w.workspaceRepo.UnsetUpToDateWithTrunkForAllInCodebase(codebaseID)
	if err != nil {
		return err
	}
	workspaces, err := w.workspaceRepo.ListByCodebaseIDs([]string{codebaseID}, false)
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
	ws, err := workspaceReader.Get(workspaceID)
	if err != nil {
		return err
	}
	t := time.Now()
	ws.UpdatedAt = &t

	// Is recalculated on next get/list
	ws.UpToDateWithTrunk = nil
	ws.HeadChangeID = nil
	ws.HeadChangeComputed = false

	if err := workspaceWriter.Update(ctx, ws); err != nil {
		return err
	}
	return nil
}
