package routes

import (
	"log"
	"net/http"

	"getsturdy.com/api/pkg/auth"
	service_jwt "getsturdy.com/api/pkg/jwt/service"
	"getsturdy.com/api/pkg/users/db"

	"github.com/gin-gonic/gin"
)

func GetSelf(repo db.Repository, jwtService *service_jwt.Service) func(c *gin.Context) {
	return func(c *gin.Context) {
		userID, err := auth.UserID(c.Request.Context())
		if err != nil {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		u, err := repo.Get(userID)
		if err != nil {
			log.Println(err)
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "failed to get user"})
			return
		}
		// Refresh the users auth cookie
		// For requests from a browser, this takes care of all of the auth renewal we need. :-)
		auth.SetAuthCookieForUser(c, u.ID, jwtService)
		c.JSON(http.StatusOK, u)
	}
}
