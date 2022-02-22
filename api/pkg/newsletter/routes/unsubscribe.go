package routes

import (
	"database/sql"
	"errors"
	"net/http"

	"getsturdy.com/api/pkg/newsletter"
	db_newsletter "getsturdy.com/api/pkg/newsletter/db"
	db_user "getsturdy.com/api/pkg/users/db"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type UnsubscribeRequest struct {
	Email string `json:"email" binding:"required"`
}

func Unsubscribe(logger *zap.Logger, userRepo db_user.Repository, notificationSettingsRepo db_newsletter.NotificationSettingsRepository) func(*gin.Context) {
	return func(c *gin.Context) {
		var req UnsubscribeRequest
		if err := c.BindJSON(&req); err != nil {
			c.AbortWithStatus(http.StatusBadRequest)
			return
		}

		logger := logger.With(zap.String("email", req.Email))

		user, err := userRepo.GetByEmail(req.Email)
		if errors.Is(err, sql.ErrNoRows) {
			// short circuit if user doesn't exist
			return
		} else if err != nil {
			logger.Error("failed to get user", zap.Error(err))
			c.AbortWithStatus(http.StatusBadRequest)
			return
		}

		settings, err := notificationSettingsRepo.GetByUser(user.ID)
		if errors.Is(err, sql.ErrNoRows) {
			if err := notificationSettingsRepo.Insert(newsletter.NotificationSettings{
				UserID:            user.ID,
				ReceiveNewsletter: false,
			}); err != nil {
				logger.Error("could not create settings", zap.Error(err))
				c.AbortWithStatus(http.StatusInternalServerError)
				return
			}
			return
		} else if err != nil {
			logger.Error("could not get settings", zap.Error(err))
			c.AbortWithStatus(http.StatusInternalServerError)
			return

		}

		settings.ReceiveNewsletter = false
		if err := notificationSettingsRepo.Update(settings); err != nil {
			logger.Error("could not update settings", zap.Error(err))
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		logger.Info("unsubscribed")

		c.Status(http.StatusOK)
	}
}
