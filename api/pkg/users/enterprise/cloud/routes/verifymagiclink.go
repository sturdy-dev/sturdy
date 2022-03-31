package routes

import (
	"errors"
	"net/http"

	"getsturdy.com/api/pkg/auth"
	service_jwt "getsturdy.com/api/pkg/jwt/service"
	"getsturdy.com/api/pkg/onetime/service"
	service_user "getsturdy.com/api/pkg/users/enterprise/cloud/service"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func VerifyMagicLink(logger *zap.Logger, userService *service_user.Service, jwtService *service_jwt.Service) gin.HandlerFunc {
	type request struct {
		Code  string `json:"code" binding:"required"`
		Email string `json:"email" binding:"required"`
	}
	return func(c *gin.Context) {
		var req request
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		user, err := userService.GetByEmail(c.Request.Context(), req.Email)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
			return
		}

		if err := userService.VerifyMagicLink(c.Request.Context(), user, req.Code); errors.Is(err, service.ErrExpired) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "code expired"})
			return
		} else if errors.Is(err, service.ErrReused) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "code already used"})
			return
		} else if errors.Is(err, service.ErrInvalid) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "code invalid"})
			return
		} else if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "something went wrong"})
			logger.Error("failed to verify magic link", zap.Error(err))
			return
		}

		if err := userService.Activate(c.Request.Context(), user); err != nil {
			logger.Error("failed to update user's status", zap.Error(err))
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		auth.SetAuthCookieForUser(c, user.ID, jwtService)
	}
}
