package routes

import (
	"net/http"
	"time"

	"getsturdy.com/api/pkg/gc/worker"
	"getsturdy.com/api/pkg/presence"
	serivce_presence "getsturdy.com/api/pkg/presence/service"
	"getsturdy.com/api/pkg/snapshots"
	db_snapshots "getsturdy.com/api/pkg/snapshots/db"
	worker_snapshots "getsturdy.com/api/pkg/snapshots/worker"
	service_suggestion "getsturdy.com/api/pkg/suggestions/service"
	db_view "getsturdy.com/api/pkg/view/db"
	"getsturdy.com/api/pkg/events"
	db_workspaces "getsturdy.com/api/pkg/workspaces/db"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type SyncTransitionsRequest struct {
	Paths      []string `json:"paths"`
	CodebaseID string   `json:"codebase_id"`
	ViewID     string   `json:"view_id"`
}

// SyncTransitions is called after mutagen has uploaded new files to the view
func SyncTransitions(
	logger *zap.Logger,
	snapshotterQueue worker_snapshots.Queue,
	viewRepo db_view.Repository,
	gcQueue *worker.Queue,
	presenceService serivce_presence.Service,
	snapshotRepo db_snapshots.Repository,
	workspaceRepo db_workspaces.WorkspaceReader,
	suggestionsService *service_suggestion.Service,
	eventSender events.EventSender,
) func(c *gin.Context) {
	return func(c *gin.Context) {
		// TODO: Authenticate internal requests

		var req SyncTransitionsRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			logger.Error("failed to parse request", zap.Error(err))
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "failed to parse or validate input"})
			return
		}

		logger.Info("sync-transition", zap.Strings("paths", req.Paths))

		if len(req.Paths) == 0 {
			c.JSON(http.StatusOK, gin.H{"status": "ok"})
			return
		}

		if err := eventSender.Codebase(req.CodebaseID, events.ViewUpdated, req.ViewID); err != nil {
			logger.Error("failed to send view updated event", zap.Error(err))
			// do not fail
		}

		// Set LastUsedAt
		view, err := viewRepo.Get(req.ViewID)
		if err != nil {
			logger.Error("could not get view", zap.Error(err))
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}
		t := time.Now()
		view.LastUsedAt = &t
		if err := viewRepo.Update(view); err != nil {
			logger.Error("could not get view", zap.Error(err))
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		if err := eventSender.Workspace(view.WorkspaceID, events.WorkspaceUpdated, view.WorkspaceID); err != nil {
			logger.Error("failed to send workspace updated event", zap.Error(err))
		}

		// Make a snapshot
		if err := snapshotterQueue.Enqueue(c, req.CodebaseID, req.ViewID, view.WorkspaceID, req.Paths, snapshots.ActionViewSync); err != nil {
			logger.Error("failed to snapshot", zap.Error(err))
			// Don't fail
		}

		if err := gcQueue.Enqueue(c, req.CodebaseID); err != nil {
			logger.Error("failed to send to gc queue", zap.Error(err))
			// Don't fail
		}

		// Create presence
		if _, err := presenceService.Record(c, view.UserID, view.WorkspaceID, presence.StateCoding); err != nil {
			logger.Error("failed to update presence", zap.Error(err))
			// Don't fail
		}

		if err := suggestionsService.RecordActivity(c, view.WorkspaceID); err != nil {
			logger.Error("failed to record suggestion activity", zap.Error(err))
			// don't fail
		}

		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	}
}
