package routes

import (
	"errors"
	"log"
	"net/http"
	"time"

	"getsturdy.com/api/pkg/auth"
	"getsturdy.com/api/pkg/jwt"
	"getsturdy.com/api/pkg/users/db"

	service_jwt "getsturdy.com/api/pkg/jwt/service"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

var (
	oneDay   = 24 * time.Hour
	oneMonth = 30 * oneDay
)

func ClientToken(db db.Repository, jwtService *service_jwt.Service) func(*gin.Context) {
	return func(c *gin.Context) {
		userID, err := auth.UserID(c.Request.Context())
		if err != nil {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		user, err := db.Get(userID)
		if err != nil {
			log.Println(err)
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "invalid input"})
			return
		}

		token, err := jwtService.IssueToken(c.Request.Context(), user.ID.String(), oneMonth, jwt.TokenTypeAuth)
		if err != nil {
			log.Println(err)
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "invalid input"})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"token": token.Token,
		})
	}
}

func RenewToken(logger *zap.Logger, db db.Repository, jwtService *service_jwt.Service) func(*gin.Context) {
	return func(c *gin.Context) {

		userID, err := auth.UserID(c.Request.Context())
		if err != nil {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		user, err := db.Get(userID)
		if err != nil {
			log.Println(err)
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "invalid input"})
			return
		}

		jwtString, err := c.Cookie("auth")
		if err != nil {
			if errors.Is(err, http.ErrNoCookie) {
				c.AbortWithStatus(http.StatusUnauthorized)
				return
			}
			c.AbortWithStatus(http.StatusBadRequest)
			return
		}

		token, err := jwtService.Verify(c.Request.Context(), jwtString, jwt.TokenTypeAuth)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "invalid input"})
			return
		}

		logger.Info("expires at", zap.Time("at", token.ExpiresAt))

		// If expires within 25 days, renew it! (original expire duration is 30 days)
		if token.ExpiresAt.Before(time.Now().Add(time.Hour * 24 * 25)) {
			newToken, err := jwtService.IssueToken(c.Request.Context(), user.ID.String(), oneMonth, jwt.TokenTypeAuth)
			if err != nil {
				logger.Error("failed to renew token for user", zap.Error(err))
				c.AbortWithStatus(http.StatusBadRequest)
				return
			}

			c.JSON(http.StatusOK, gin.H{
				"token":   newToken.Token,
				"has_new": true,
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"token":   "",
			"has_new": false,
		})
	}
}
