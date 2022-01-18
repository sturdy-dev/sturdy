package routes

import (
	"mash/pkg/auth"
	"mash/pkg/change/vcs"
	"mash/pkg/snapshots"
	"mash/pkg/view/meta"
	"mash/vcs/executor"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"mash/pkg/codebase/access"
	db_codebase "mash/pkg/codebase/db"
	"mash/pkg/view/db"
)

type IgnoreFileRequest struct {
	Path string `json:"path" binding:"required"`
}

func IgnoreFile(logger *zap.Logger, viewRepo db.Repository, codebaseUserRepo db_codebase.CodebaseUserRepository, executorProvider executor.Provider, viewUpdatedFunc meta.ViewUpdatedFunc) func(c *gin.Context) {
	return func(c *gin.Context) {
		userID, err := auth.UserID(c.Request.Context())
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, err)
			return
		}

		var req IgnoreFileRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			logger.Warn("parse request failed", zap.Error(err))
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "failed to parse or validate input"})
			return
		}

		viewID := c.Param("viewID")
		view, err := viewRepo.Get(viewID)
		if err != nil {
			logger.Error("failed to get repo", zap.Error(err))
			c.AbortWithStatusJSON(http.StatusBadRequest, err)
			return
		}

		// gitignores can not be ignored
		if strings.Contains(req.Path, ".gitignore") {
			c.Status(http.StatusOK)
			return
		}

		if !access.UserHasAccessToCodebase(codebaseUserRepo, userID, view.CodebaseID) {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		// Update .gitignore file
		err = vcs.AddToGitignore(executorProvider, view.CodebaseID, view.ID, req.Path)
		if err != nil {
			logger.Error("failed to add ignore", zap.Error(err))
			c.AbortWithStatusJSON(http.StatusInternalServerError, err)
			return
		}

		err = viewUpdatedFunc(c, view, snapshots.ActionFileIgnore)
		if err != nil {
			logger.Error("failed to mark as updated", zap.Error(err))
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		c.Status(http.StatusOK)
	}
}
