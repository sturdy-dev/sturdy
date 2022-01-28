package routes

import (
	"errors"
	"net/http"
	"strings"

	"getsturdy.com/api/pkg/analytics"
	"getsturdy.com/api/pkg/auth"
	service_jwt "getsturdy.com/api/pkg/jwt/service"
	service_user "getsturdy.com/api/pkg/user/service"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func Signup(logger *zap.Logger, userService service_user.Service, jwtService *service_jwt.Service, analyticsClient analytics.Client) func(c *gin.Context) {
	type request struct {
		Name     string `json:"name" binding:"required"`
		Email    string `json:"email" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	return func(c *gin.Context) {
		var req request
		if err := c.ShouldBindJSON(&req); err != nil {
			logger.Warn("failed to bind input", zap.Error(err))
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "invalid input"})
			return
		}

		req.Email = strings.TrimSpace(req.Email)
		req.Name = strings.TrimSpace(req.Name)
		req.Email = strings.TrimSpace(req.Email)

		newUser, err := userService.CreateWithPassword(c.Request.Context(), req.Name, req.Password, req.Email)
		if errors.Is(err, service_user.ErrExists) {
			c.AbortWithStatusJSON(http.StatusConflict, gin.H{"error": "User already exists"})
			return
		} else if err != nil {
			logger.Error("failed to create user", zap.Error(err))
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		if err := auth.SetAuthCookieForUser(c, newUser.ID, jwtService); err != nil {
			logger.Error("failed to set auth cookie", zap.Error(err))
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		if err := analyticsClient.Enqueue(analytics.Capture{
			DistinctId: newUser.ID,
			Event:      "logged in",
			Properties: analytics.NewProperties().
				Set("type", "password"),
		}); err != nil {
			logger.Error("send to analytics failed", zap.Error(err))
		}

		// Send the user object in the response
		c.JSON(http.StatusOK, newUser)
	}
}
