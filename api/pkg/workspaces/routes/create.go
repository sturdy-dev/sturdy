package routes

import (
	"net/http"
	"strings"

	"getsturdy.com/api/pkg/auth"
	service_workspace "getsturdy.com/api/pkg/workspaces/service"

	"go.uber.org/zap"

	"github.com/gin-gonic/gin"

	"getsturdy.com/api/pkg/codebase/access"
	db_codebase "getsturdy.com/api/pkg/codebase/db"
)

type CreateRequest struct {
	CodebaseID string `json:"codebase_id" binding:"required"`
	Name       string `json:"name"`

	// ChangeID is a commit checksum
	ChangeID string `json:"change_id"` // change_id and revert_change_id are mutually exclusive

	// RevertChangeID is a commit checksum
	RevertChangeID string `json:"revert_change_id"` //
}

func Create(logger *zap.Logger, workspaceService service_workspace.Service, codebaseUserRepo db_codebase.CodebaseUserRepository) func(c *gin.Context) {
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
		req.CodebaseID = strings.TrimSpace(req.CodebaseID)

		if !access.UserHasAccessToCodebase(codebaseUserRepo, userID, req.CodebaseID) {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		ws, err := workspaceService.Create(c.Request.Context(), service_workspace.CreateWorkspaceRequest{
			UserID:         userID,
			CodebaseID:     req.CodebaseID,
			Name:           req.Name,
			RevertCommitID: req.RevertChangeID,
			CommitID:       req.ChangeID,
		})

		if err != nil {
			logger.Error("failed to create workspace", zap.Error(err))
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		c.JSON(http.StatusOK, ws)
	}
}
