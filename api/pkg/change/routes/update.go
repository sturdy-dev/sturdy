package routes

import (
	"database/sql"
	"errors"
	"getsturdy.com/api/pkg/analytics"
	"getsturdy.com/api/pkg/auth"
	"getsturdy.com/api/pkg/change"
	"getsturdy.com/api/pkg/change/message"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	changeDB "getsturdy.com/api/pkg/change/db"
	"getsturdy.com/api/pkg/codebase/access"
	codebaseDB "getsturdy.com/api/pkg/codebase/db"
)

type UpdateRequest struct {
	UpdatedDescription string `json:"updated_description" binding:"required"`
}

func Update(logger *zap.Logger, codebaseUserRepo codebaseDB.CodebaseUserRepository, analyticsClient analytics.Client, changeRepo changeDB.Repository) func(c *gin.Context) {
	return func(c *gin.Context) {
		userID, err := auth.UserID(c.Request.Context())
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

		changeID := change.ID(c.Param("id"))

		// Get the change (it might exist in the db)
		ch, err := changeRepo.Get(changeID)
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
			cleanCommitMessageTitle := strings.Split(cleanCommitMessage, "\n")[0]
			ch.Title = &cleanCommitMessageTitle
		}

		err = changeRepo.Update(ch)
		if err != nil {
			logger.Error("failed to update change", zap.Error(err))
			c.AbortWithStatusJSON(http.StatusBadRequest, err)
			return
		}

		err = analyticsClient.Enqueue(analytics.Capture{
			DistinctId: userID,
			Event:      "updated change",
			Properties: analytics.NewProperties().
				Set("commit_id", changeID).
				Set("codebase_id", ch.CodebaseID),
		})
		if err != nil {
			logger.Error("analytics failed", zap.Error(err))
		}

		// TODO: Migrate this to GraphQL, it's a temporary hack for now
		c.JSON(http.StatusOK, gin.H{
			"title":       ch.Title,
			"description": ch.UpdatedDescription,
		})
	}
}
