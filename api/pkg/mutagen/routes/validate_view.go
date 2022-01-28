package routes

import (
	"getsturdy.com/api/pkg/analytics"
	"getsturdy.com/api/pkg/events"
	"net/http"
	"time"

	"go.uber.org/zap"

	"github.com/gin-gonic/gin"

	db_view "getsturdy.com/api/pkg/view/db"
)

type ValidateViewRequest struct {
	ViewID          string `json:"view_id" binding:"required"`
	CodebaseID      string `json:"codebase_id" binding:"required"`
	UserID          string `json:"user_id" binding:"required"`
	IsNewConnection bool   `json:"is_new_connection"`
}

func ValidateView(logger *zap.Logger, viewRepo db_view.Repository, analyticsClient analytics.Client, eventsSender events.EventSender) func(c *gin.Context) {
	return func(c *gin.Context) {
		var req ValidateViewRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			logger.Error("failed to parse request", zap.Error(err))
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "failed to parse or validate input"})
			return
		}

		logger.Info("validate", zap.Any("req", req))

		viewObj, err := viewRepo.Get(req.ViewID)
		if err != nil {
			logger.Warn("view not found", zap.Error(err))
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": "view not found"})
			return
		}

		if viewObj.CodebaseID != req.CodebaseID {
			logger.Warn("codebase did not match", zap.Error(err))
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "codebase not found"})
			return
		}

		if viewObj.UserID != req.UserID {
			logger.Warn("user did not match", zap.Error(err))
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "user not found"})
			return
		}

		// Set LastUsedAt
		t := time.Now()
		viewObj.LastUsedAt = &t
		err = viewRepo.Update(viewObj)
		if err != nil {
			logger.Error("could not update view", zap.Error(err))
			c.AbortWithStatus(http.StatusInternalServerError)
		}

		if err := eventsSender.Codebase(viewObj.CodebaseID, events.ViewUpdated, viewObj.ID); err != nil {
			logger.Error("could not send event", zap.Error(err))
			// do not fail
		}

		if req.IsNewConnection {
			err = analyticsClient.Enqueue(analytics.Capture{
				DistinctId: req.UserID,
				Event:      "mutagen connection to view",
				Properties: analytics.NewProperties().
					Set("codebase_id", req.CodebaseID).
					Set("view_id", req.ViewID).
					Set("is_new_connection", req.IsNewConnection), // This is always true...
			})
			if err != nil {
				logger.Error("analytics failed", zap.Error(err))
			}
		}

		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	}
}
