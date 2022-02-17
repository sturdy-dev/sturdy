package routes

import (
	"net/http"
	"time"

	"getsturdy.com/api/pkg/analytics"
	service_analytics "getsturdy.com/api/pkg/analytics/service"
	"getsturdy.com/api/pkg/events"
	db_view "getsturdy.com/api/pkg/view/db"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type ValidateViewRequest struct {
	ViewID          string `json:"view_id" binding:"required"`
	CodebaseID      string `json:"codebase_id" binding:"required"`
	UserID          string `json:"user_id" binding:"required"`
	IsNewConnection bool   `json:"is_new_connection"`
}

func ValidateView(logger *zap.Logger, viewRepo db_view.Repository, analyticsService *service_analytics.Service, eventsSender events.EventSender) func(c *gin.Context) {
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
			analyticsService.Capture(c.Request.Context(), "mutagen connection to view",
				analytics.DistinctID(req.UserID),
				analytics.CodebaseID(req.CodebaseID),
				analytics.Property("view_id", req.ViewID),
				analytics.Property("is_new_connection", req.IsNewConnection), // This is alway true...
			)
		}

		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	}
}
