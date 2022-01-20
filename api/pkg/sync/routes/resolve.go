package routes

import (
	service_sync "getsturdy.com/api/pkg/sync/service"
	vcsvcs "getsturdy.com/api/vcs"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type InitSyncRequest struct {
	WorkspaceID string `json:"workspace_id" binding:"required"`
}

type ResolveRequest struct {
	Files []ResolveFileRequest `json:"files" binding:"required"`
}

type ResolveFileRequest struct {
	FilePath string `json:"file_path" binding:"required"`
	Version  string `json:"version" binding:"required"`
}

func ResolveV2(
	logger *zap.Logger,
	syncService *service_sync.Service,
) func(c *gin.Context) {
	return func(c *gin.Context) {
		var req ResolveRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			logger.Error("failed to parse request", zap.Error(err))
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "failed to parse or validate input"})
			return
		}

		viewID := c.Param("viewID")

		var resolves []vcsvcs.SturdyRebaseResolve
		for _, rf := range req.Files {
			resolves = append(resolves, vcsvcs.SturdyRebaseResolve{Path: rf.FilePath, Version: rf.Version})
		}

		if status, err := syncService.Resolve(viewID, resolves); err != nil {
			logger.Error("failed to sync on trunk", zap.Error(err))
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		} else {
			c.JSON(http.StatusOK, status)
			return
		}
	}
}
