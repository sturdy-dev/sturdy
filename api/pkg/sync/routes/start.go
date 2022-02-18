package routes

import (
	"net/http"

	"go.uber.org/zap"

	"github.com/gin-gonic/gin"

	service_sync "getsturdy.com/api/pkg/sync/service"
)

// Auth happens elsewhere
func StartV2(
	logger *zap.Logger,
	syncService *service_sync.Service,
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

		if status, err := syncService.OnTrunk(c.Request.Context(), viewID); err != nil {
			logger.Error("failed to sync on trunk", zap.Error(err))
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		} else {
			c.JSON(http.StatusOK, status)
			return
		}
	}
}
