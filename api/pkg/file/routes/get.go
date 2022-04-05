package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	service_auth "getsturdy.com/api/pkg/auth/service"
	service_file "getsturdy.com/api/pkg/file/service"
	service_workspace "getsturdy.com/api/pkg/workspaces/service"
)

type GetFileRoute func(*gin.Context)

func NewGetFileRoute(
	workspaceService service_workspace.Service,
	authService *service_auth.Service,
	fileService *service_file.Service,
	logger *zap.Logger,
) GetFileRoute {
	return func(c *gin.Context) {
		ctx := c.Request.Context()

		path := c.Query("path")
		workspaceID := c.Query("workspace_id")
		isNew := c.Query("is_new")

		ws, err := workspaceService.GetByID(ctx, workspaceID)
		if err != nil {
			logger.Error("could not get workspace", zap.Error(err))
			c.AbortWithStatus(http.StatusNotFound)
			return
		}

		if err := authService.CanRead(ctx, ws); err != nil {
			c.AbortWithStatus(http.StatusNotFound)
			return
		}

		data, err := fileService.ReadWorkspaceFile(ctx, ws, path, isNew == "1")
		if err != nil {
			logger.Error("could not get file", zap.Error(err))
			c.AbortWithStatus(http.StatusNotFound)
			return
		}

		c.Status(http.StatusOK)
		c.Writer.Write(data)
	}
}
