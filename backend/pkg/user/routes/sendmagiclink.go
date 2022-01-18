package routes

import (
	service_user "mash/pkg/user/service"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func SendMagicLink(logger *zap.Logger, userService *service_user.Service) gin.HandlerFunc {
	type request struct {
		Name  string `json:"name"`
		Email string `json:"email" binding:"required"`
	}

	return func(c *gin.Context) {
		var req request
		if err := c.ShouldBindJSON(&req); err != nil {
			logger.Warn("failed to bind input", zap.Error(err))
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "invalid input"})
			return
		}

		// If the user already exists, an OTP email will be sent to that user
		// If the user doesn't exist, the user will be created (maybe without a name)

		req.Email = strings.TrimSpace(req.Email)
		req.Name = strings.TrimSpace(req.Name)

		newUser, err := userService.Create(c.Request.Context(), req.Name, req.Email)
		if err != nil {
			logger.Warn("failed to create user", zap.Error(err))
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "failed to create user"})
			return
		}

		c.JSON(http.StatusOK, newUser)
	}
}
