package routes

import (
	"net/http"
	"strings"

	"getsturdy.com/api/pkg/analytics"
	"getsturdy.com/api/pkg/auth"
	service_jwt "getsturdy.com/api/pkg/jwt/service"
	"getsturdy.com/api/pkg/users/db"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

func Login(logger *zap.Logger, repo db.Repository, analyticsClient analytics.Client, jwtService *service_jwt.Service) func(c *gin.Context) {
	type request struct {
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

		// Get user by email
		getUser, err := repo.GetByEmail(req.Email)
		if err != nil {
			logger.Warn("failed to get user", zap.Error(err))
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "invalid password"})
			return
		}

		// Compare the input password against our hashed version
		if err := bcrypt.CompareHashAndPassword(
			[]byte(getUser.PasswordHash),
			[]byte(req.Password),
		); err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "invalid password"})
			return
		}

		if err := auth.SetAuthCookieForUser(c, getUser.ID, jwtService); err != nil {
			logger.Error("failed to set auth cookie", zap.Error(err))
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		// Identify to analytics
		if err := analyticsClient.Enqueue(analytics.Identify{
			DistinctId: getUser.ID,
			Properties: analytics.NewProperties().
				Set("name", getUser.Name).
				Set("email", getUser.Email),
		}); err != nil {
			logger.Error("send to analytics failed", zap.Error(err))
		}

		if err := analyticsClient.Enqueue(analytics.Capture{
			DistinctId: getUser.ID,
			Event:      "logged in",
			Properties: analytics.NewProperties().
				Set("type", "password"),
		}); err != nil {
			logger.Error("send to analytics failed", zap.Error(err))
		}

		// Send the user object in the response
		c.JSON(http.StatusOK, getUser)
	}
}
