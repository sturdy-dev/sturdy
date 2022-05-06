package routes

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"getsturdy.com/api/pkg/codebases"
	eventsv2 "getsturdy.com/api/pkg/events/v2"
	"getsturdy.com/api/pkg/gc/worker"
	"getsturdy.com/api/pkg/presence"
	serivce_presence "getsturdy.com/api/pkg/presence/service"
	"getsturdy.com/api/pkg/snapshots"
	worker_snapshots "getsturdy.com/api/pkg/snapshots/worker"
	service_suggestion "getsturdy.com/api/pkg/suggestions/service"
	db_view "getsturdy.com/api/pkg/views/db"
)

type SyncTransitionsRequest struct {
	Paths      []string     `json:"paths"`
	CodebaseID codebases.ID `json:"codebase_id"`
	ViewID     string       `json:"view_id"`
}

// SyncTransitions is called after mutagen has uploaded new files to the view
func SyncTransitions(
	logger *zap.Logger,
	snapshotterQueue worker_snapshots.Queue,
	viewRepo db_view.Repository,
	gcQueue *worker.Queue,
	presenceService serivce_presence.Service,
	suggestionsService *service_suggestion.Service,
	eventSender *eventsv2.Publisher,
) func(c *gin.Context) {
	return func(c *gin.Context) {
		// TODO: Authenticate internal requests

		ctx := c.Request.Context()

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

		view, err := viewRepo.Get(req.ViewID)
		if err != nil {
			logger.Error("could not get view", zap.Error(err))
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		if err := eventSender.ViewUpdated(ctx, eventsv2.Codebase(req.CodebaseID), view); err != nil {
			logger.Error("failed to send view updated event", zap.Error(err))
			// do not fail
		}

		// Set LastUsedAt
		t := time.Now()
		view.LastUsedAt = &t
		if err := viewRepo.Update(view); err != nil {
			logger.Error("could not get view", zap.Error(err))
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		// Make a snapshot
		if err := snapshotterQueue.Enqueue(c, req.CodebaseID, req.ViewID, view.WorkspaceID, view.UserID, snapshots.ActionViewSync); err != nil {
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
