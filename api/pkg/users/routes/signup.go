package routes

import (
	"errors"
	"net/http"
	"strings"

	"getsturdy.com/api/pkg/analytics"
	service_analytics "getsturdy.com/api/pkg/analytics/service"
	"getsturdy.com/api/pkg/auth"
	service_jwt "getsturdy.com/api/pkg/jwt/service"
	service_user "getsturdy.com/api/pkg/users/service"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func Signup(
	logger *zap.Logger,
	userService service_user.Service,
	jwtService *service_jwt.Service,
	analyticsService *service_analytics.Service,
) func(c *gin.Context) {
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
			c.AbortWithStatusJSON(http.StatusConflict, gin.H{"error": "A user with this email already exists, did you mean to login?"})
			return
		} else if errors.Is(err, service_user.ErrExceeded) {
			c.AbortWithStatusJSON(http.StatusConflict, gin.H{"error": "This service is exceeding the number of users allowed. Please contact the server administrator or email support@getsturdy.com"})
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

		ctx := c.Request.Context()
		analyticsService.IdentifyUser(ctx, newUser)
		analyticsService.CaptureUser(newUser.ID, "logged in", analytics.Property("type", "password"))

		// Send the user object in the response
		c.JSON(http.StatusOK, newUser)
	}
}
