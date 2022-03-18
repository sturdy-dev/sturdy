package routes

import (
	"net/http"
	"strings"

	"getsturdy.com/api/pkg/auth"
	"getsturdy.com/api/pkg/codebases"
	service_workspace "getsturdy.com/api/pkg/workspaces/service"

	"go.uber.org/zap"

	"github.com/gin-gonic/gin"

	"getsturdy.com/api/pkg/codebases/access"
	db_codebases "getsturdy.com/api/pkg/codebases/db"
)

type CreateRequest struct {
	CodebaseID codebases.ID `json:"codebase_id" binding:"required"`
}

func Create(logger *zap.Logger, workspaceService service_workspace.Service, codebaseUserRepo db_codebases.CodebaseUserRepository) func(c *gin.Context) {
	return func(c *gin.Context) {
		userID, err := auth.UserID(c.Request.Context())
		if err != nil {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		var req CreateRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			logger.Warn("failed to parse input", zap.Error(err))
			c.AbortWithStatus(http.StatusBadRequest)
			return
		}
		req.CodebaseID = codebases.ID(strings.TrimSpace(req.CodebaseID.String()))

		if !access.UserHasAccessToCodebase(codebaseUserRepo, userID, req.CodebaseID) {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		ws, err := workspaceService.Create(c.Request.Context(), service_workspace.CreateWorkspaceRequest{
			UserID:     userID,
			CodebaseID: req.CodebaseID,
		})
		if err != nil {
			logger.Error("failed to create workspace", zap.Error(err))
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		c.JSON(http.StatusOK, ws)
	}
}
