package routes

import (
	"database/sql"
	"errors"
	"net/http"
	"time"

	service_auth "getsturdy.com/api/pkg/auth/service"
	"getsturdy.com/api/pkg/codebase"
	"getsturdy.com/api/pkg/codebase/db"
	service_codebase "getsturdy.com/api/pkg/codebase/service"
	userDB "getsturdy.com/api/pkg/user/db"
	"getsturdy.com/api/pkg/events"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type InviteUserRequest struct {
	UserEmail string `json:"user_email" binding:"required"`
}

func Invite(userRepo userDB.Repository, codeBaseUserRepo db.CodebaseUserRepository, codebaseService *service_codebase.Service, authService *service_auth.Service, eventsSender events.EventSender, logger *zap.Logger) func(c *gin.Context) {
	return func(c *gin.Context) {
		codebaseID := c.Param("id")

		cb, err := codebaseService.GetByID(c.Request.Context(), codebaseID)
		if err != nil && errors.Is(err, sql.ErrNoRows) {
			c.AbortWithStatus(http.StatusNotFound)
			return
		} else if err != nil {
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		if err := authService.CanWrite(c.Request.Context(), cb); err != nil {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		var request InviteUserRequest
		if err := c.BindJSON(&request); err != nil {
			c.AbortWithStatus(http.StatusBadRequest)
			return
		}

		inviteUser, err := userRepo.GetByEmail(request.UserEmail)
		if err != nil {
			logger.Error("failed to invite user", zap.Error(err))
			c.AbortWithStatus(http.StatusBadRequest)
			return
		}

		// Check that the user isn't already a member
		_, err = codeBaseUserRepo.GetByUserAndCodebase(inviteUser.ID, codebaseID)
		if err == nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "invalid request - user is already a member"})
			return
		}

		t := time.Now()
		err = codeBaseUserRepo.Create(codebase.CodebaseUser{
			ID:         uuid.New().String(),
			UserID:     inviteUser.ID,
			CodebaseID: codebaseID,
			CreatedAt:  &t,
		})
		if err != nil {
			logger.Error("failed to invite user", zap.Error(err))
			c.AbortWithStatus(http.StatusBadRequest)
			return
		}

		// Send events
		if err := eventsSender.Codebase(codebaseID, events.CodebaseUpdated, codebaseID); err != nil {
			logger.Error("failed to send events", zap.Error(err))
		}

		c.JSON(http.StatusOK, inviteUser)
	}
}
