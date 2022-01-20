package routes

import (
	"database/sql"
	"errors"
	"fmt"
	"getsturdy.com/api/pkg/author"
	"getsturdy.com/api/pkg/jsontime"
	db_user "getsturdy.com/api/pkg/user/db"
	"getsturdy.com/api/pkg/view"
	"getsturdy.com/api/pkg/view/db"
	db_workspace "getsturdy.com/api/pkg/workspace/db"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func Get(repo db.Repository, workspaceReader db_workspace.WorkspaceReader, userRepo db_user.Repository, logger *zap.Logger) func(c *gin.Context) {
	return func(c *gin.Context) {
		viewID := c.Param("viewID")

		viewObj, err := repo.Get(viewID)
		if err != nil {
			logger.Error("failed to get view", zap.Error(err))
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "failed to get view"})
			return
		}

		res, err := addMeta(toViewJson(viewObj), workspaceReader, userRepo)
		if err != nil {
			logger.Error("failed to get view meta", zap.Error(err))
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "failed to get view"})
			return
		}

		c.JSON(http.StatusOK, res)
	}
}

func toViewJson(in *view.View) view.ViewJSON {
	var lastUsedAt jsontime.Time
	if in.LastUsedAt != nil {
		lastUsedAt = jsontime.Time(*in.LastUsedAt)
	}

	return view.ViewJSON{
		View:       *in,
		LastUsedAt: lastUsedAt,
	}
}

func addMeta(v view.ViewJSON, workspaceReader db_workspace.WorkspaceReader, userRepo db_user.Repository) (view.ViewWithMetadataJSON, error) {
	user, err := author.GetAuthor(v.UserID, userRepo)
	if err != nil {
		return view.ViewWithMetadataJSON{}, fmt.Errorf("failed to get user metadata: %w", err)
	}

	res := view.ViewWithMetadataJSON{
		ViewJSON: v,
		User:     user,
	}

	ws, err := workspaceReader.Get(v.WorkspaceID)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return view.ViewWithMetadataJSON{}, fmt.Errorf("failed to get workspace metadata: %w", err)
	}
	if err == nil {
		res.ViewWorkspaceMeta = view.ViewWorkspaceMeta{
			ID:   ws.ID,
			Name: ws.Name,
		}
	}

	return res, nil
}
