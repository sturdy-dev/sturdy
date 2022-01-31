package routes

import (
	"database/sql"
	"errors"
	"getsturdy.com/api/pkg/newsletter"
	db_newsletter "getsturdy.com/api/pkg/newsletter/db"
	db_user "getsturdy.com/api/pkg/users/db"
	"net/http"

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
		if err != nil {
			logger.Error("could not get user", zap.Error(err))
			c.AbortWithStatus(http.StatusBadRequest)
			return
		}

		settings, err := notificationSettingsRepo.GetByUser(user.ID)
		if errors.Is(err, sql.ErrNoRows) {
			err = notificationSettingsRepo.Insert(newsletter.NotificationSettings{
				UserID:            user.ID,
				ReceiveNewsletter: false,
			})
			if err != nil {
				logger.Error("could not create settings", zap.Error(err))
				c.AbortWithStatus(http.StatusInternalServerError)
				return
			}
			c.Status(http.StatusOK)
			return
		} else if err != nil {
			logger.Error("could not get settings", zap.Error(err))
			c.AbortWithStatus(http.StatusInternalServerError)
			return

		}

		settings.ReceiveNewsletter = false
		err = notificationSettingsRepo.Update(settings)
		if err != nil {
			logger.Error("could not update settings", zap.Error(err))
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		logger.Info("unsubscribed")

		c.Status(http.StatusOK)
	}
}
