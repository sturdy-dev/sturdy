package routes

import (
	"database/sql"
	"errors"
	"net/http"

	"getsturdy.com/api/pkg/analytics"
	service_analytics "getsturdy.com/api/pkg/analytics/service"
	"getsturdy.com/api/pkg/auth"
	"getsturdy.com/api/pkg/changes"
	changeDB "getsturdy.com/api/pkg/changes/db"
	"getsturdy.com/api/pkg/changes/message"
	"getsturdy.com/api/pkg/codebases/access"
	codebaseDB "getsturdy.com/api/pkg/codebases/db"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type UpdateRequest struct {
	UpdatedDescription string `json:"updated_description" binding:"required"`
}

func Update(
	logger *zap.Logger,
	codebaseUserRepo codebaseDB.CodebaseUserRepository,
	analyticsService *service_analytics.Service,
	changeRepo changeDB.Repository,
) func(c *gin.Context) {
	return func(c *gin.Context) {
		ctx := c.Request.Context()

		userID, err := auth.UserID(ctx)
		if err != nil {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		var req UpdateRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			logger.Warn("parse request failed", zap.Error(err))
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "failed to parse or validate input"})
			return
		}

		changeID := changes.ID(c.Param("id"))

		// Get the change (it might exist in the db)
		ch, err := changeRepo.Get(ctx, changeID)
		if err != nil && errors.Is(err, sql.ErrNoRows) {
			c.AbortWithStatus(http.StatusNotFound)
			return
		} else if err != nil {
			logger.Error("failed to get change", zap.Error(err))
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		if !access.UserHasAccessToCodebase(codebaseUserRepo, userID, ch.CodebaseID) {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		// Update the change
		if len(req.UpdatedDescription) > 0 {
			ch.UpdatedDescription = req.UpdatedDescription

			// Use the first line as the title
			cleanCommitMessage := message.CommitMessage(req.UpdatedDescription)
			title := message.Title(cleanCommitMessage)
			ch.Title = &title
		}

		err = changeRepo.Update(ctx, *ch)
		if err != nil {
			logger.Error("failed to update change", zap.Error(err))
			c.AbortWithStatusJSON(http.StatusBadRequest, err)
			return
		}

		analyticsService.Capture(c.Request.Context(), "updated change",
			analytics.CodebaseID(ch.CodebaseID),
			analytics.Property("commit_id", ch.ID),
		)

		// TODO: Migrate this to GraphQL, it's a temporary hack for now
		c.JSON(http.StatusOK, gin.H{
			"title":       ch.Title,
			"description": ch.UpdatedDescription,
		})
	}
}
