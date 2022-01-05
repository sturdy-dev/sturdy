package routes

import (
	"database/sql"
	"errors"
	"mash/pkg/auth"
	"mash/pkg/view/events"
	"net/http"
	"time"

	"github.com/google/uuid"

	"mash/pkg/codebase"

	"go.uber.org/zap"

	"mash/pkg/codebase/db"

	"github.com/gin-gonic/gin"
)

func JoinGetCodebase(logger *zap.Logger, repo db.CodebaseRepository) func(c *gin.Context) {
	return func(c *gin.Context) {
		code := c.Param("code")
		if len(code) == 0 {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "unable to get codebase"})
			return
		}

		codebase, err := repo.GetByInviteCode(code)
		if err != nil {
			if !errors.Is(err, sql.ErrNoRows) {
				logger.Error("could not get codebase with join code", zap.Error(err))
			}
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "unable to get codebase"})
			return
		}

		c.JSON(http.StatusOK, codebase)
	}
}

func JoinCodebase(logger *zap.Logger, repo db.CodebaseRepository, codeBaseUserRepo db.CodebaseUserRepository, eventSender events.EventSender) func(c *gin.Context) {
	return func(c *gin.Context) {
		code := c.Param("code")
		if len(code) == 0 {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "unable to get codebase"})
			return
		}

		cb, err := repo.GetByInviteCode(code)
		if err != nil {
			if !errors.Is(err, sql.ErrNoRows) {
				logger.Error("could not get codebase with join code", zap.Error(err))
			}
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "unable to get codebase"})
			return
		}

		userID, err := auth.UserID(c.Request.Context())
		if err != nil {
			c.AbortWithStatus(http.StatusBadRequest)
			return
		}

		// Check that the user isn't already a member
		_, err = codeBaseUserRepo.GetByUserAndCodebase(userID, cb.ID)
		if err == nil {
			c.JSON(http.StatusOK, cb)
			return
		}

		t := time.Now()
		err = codeBaseUserRepo.Create(codebase.CodebaseUser{
			ID:         uuid.New().String(),
			UserID:     userID,
			CodebaseID: cb.ID,
			CreatedAt:  &t,
		})
		if err != nil {
			logger.Error("failed to invite user", zap.Error(err))
			c.AbortWithStatus(http.StatusBadRequest)
			return
		}

		// Send events
		if err := eventSender.Codebase(cb.ID, events.CodebaseUpdated, cb.ID); err != nil {
			logger.Error("failed to send events", zap.Error(err))
		}

		c.JSON(http.StatusOK, cb)
	}
}
