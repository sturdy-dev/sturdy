package routes

import (
	"net/http"
	"strings"

	"getsturdy.com/api/pkg/analytics"
	service_analytics "getsturdy.com/api/pkg/analytics/service"
	"getsturdy.com/api/pkg/auth"
	service_jwt "getsturdy.com/api/pkg/jwt/service"
	service_users "getsturdy.com/api/pkg/users/service"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

func Login(logger *zap.Logger, userService service_users.Service, analyticsService *service_analytics.Service, jwtService *service_jwt.Service) func(c *gin.Context) {
	type request struct {
		Email    string `json:"email" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	return func(c *gin.Context) {
		var req request
		if err := c.ShouldBindJSON(&req); err != nil {
			logger.Warn("failed to bind input", zap.Error(err))
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Something went wrong, please check the input and try again"})
			return
		}

		req.Email = strings.TrimSpace(req.Email)

		// Get user by email
		getUser, err := userService.GetByEmail(c.Request.Context(), req.Email)
		if err != nil {
			logger.Warn("failed to get user", zap.Error(err))
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid email or password, please check the input and try again"})
			return
		}

		// Compare the input password against our hashed version
		if err := bcrypt.CompareHashAndPassword(
			[]byte(getUser.PasswordHash),
			[]byte(req.Password),
		); err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid email or password, please check the input and try again"})
			return
		}

		if err := auth.SetAuthCookieForUser(c, getUser.ID, jwtService); err != nil {
			logger.Error("failed to set auth cookie", zap.Error(err))
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		if err := userService.Activate(c.Request.Context(), getUser); err != nil {
			logger.Error("failed to activate user", zap.Error(err))
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		ctx := c.Request.Context()
		analyticsService.IdentifyUser(ctx, getUser)
		analyticsService.Capture(ctx, "logged in", analytics.Property("type", "password"))

		// Send the user object in the response
		c.JSON(http.StatusOK, getUser)
	}
}
