package routes

import (
	"database/sql"
	"errors"
	"net/http"

	"go.uber.org/zap"

	"github.com/gin-gonic/gin"

	service_sync "getsturdy.com/api/pkg/sync/service"
	service_workspace "getsturdy.com/api/pkg/workspaces/service"
)

// Auth happens elsewhere
func StartV2(
	logger *zap.Logger,
	syncService *service_sync.Service,
	workspaceService service_workspace.Service,
) func(c *gin.Context) {
	return func(c *gin.Context) {
		var req InitSyncRequest
		viewID := c.Param("viewID")
		if err := c.ShouldBindJSON(&req); err != nil {
			logger.Error("failed to parse input", zap.Error(err))
			c.AbortWithStatus(http.StatusBadRequest)
			return
		}

		logger := logger.With(zap.String("view_id", viewID))

		workspace, err := workspaceService.GetByViewID(c.Request.Context(), viewID)
		if errors.Is(err, sql.ErrNoRows) {
			c.AbortWithStatus(http.StatusNotFound)
			return
		} else if err != nil {
			logger.Error("failed to get view", zap.Error(err))
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		if status, err := syncService.OnTrunk(c.Request.Context(), workspace); err != nil {
			logger.Error("failed to sync on trunk", zap.Error(err))
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		} else {
			c.JSON(http.StatusOK, status)
			return
		}
	}
}
