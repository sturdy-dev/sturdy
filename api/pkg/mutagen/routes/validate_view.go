package routes

import (
	"net/http"
	"time"

	"getsturdy.com/api/pkg/analytics"
	service_analytics "getsturdy.com/api/pkg/analytics/service"
	"getsturdy.com/api/pkg/codebases"
	eventsv2 "getsturdy.com/api/pkg/events/v2"
	"getsturdy.com/api/pkg/users"
	db_view "getsturdy.com/api/pkg/view/db"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type ValidateViewRequest struct {
	ViewID          string       `json:"view_id" binding:"required"`
	CodebaseID      codebases.ID `json:"codebase_id" binding:"required"`
	UserID          users.ID     `json:"user_id" binding:"required"`
	IsNewConnection bool         `json:"is_new_connection"`
}

func ValidateView(logger *zap.Logger, viewRepo db_view.Repository, analyticsService *service_analytics.Service, eventsSender *eventsv2.Publisher) func(c *gin.Context) {
	return func(c *gin.Context) {
		ctx := c.Request.Context()

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

		if err := eventsSender.ViewUpdated(ctx, eventsv2.Codebase(viewObj.CodebaseID), viewObj); err != nil {
			logger.Error("could not send event", zap.Error(err))
			// do not fail
		}

		if req.IsNewConnection {
			analyticsService.Capture(c.Request.Context(), "mutagen connection to view",
				analytics.UserID(req.UserID),
				analytics.CodebaseID(req.CodebaseID),
				analytics.Property("view_id", req.ViewID),
				analytics.Property("is_new_connection", req.IsNewConnection), // This is alway true...
			)
		}

		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	}
}
