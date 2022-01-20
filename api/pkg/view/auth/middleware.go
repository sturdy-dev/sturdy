package auth

import (
	"database/sql"
	"errors"
	"net/http"

	service_auth "getsturdy.com/api/pkg/auth/service"
	db_view "getsturdy.com/api/pkg/view/db"

	"github.com/gin-gonic/gin"
)

func ValidateViewAccessMiddleware(authService *service_auth.Service, viewRepo db_view.Repository) func(*gin.Context) {
	return func(c *gin.Context) {
		id := c.Param("viewID")

		view, err := viewRepo.Get(id)
		if err != nil && errors.Is(err, sql.ErrNoRows) {
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		} else if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "view not found"})
			return
		}

		switch c.Request.Method {
		case http.MethodGet:
			if err := authService.CanRead(c.Request.Context(), view); err != nil {
				c.AbortWithStatus(http.StatusUnauthorized)
				return
			}
		case http.MethodPost:
			if err := authService.CanWrite(c.Request.Context(), view); err != nil {
				c.AbortWithStatus(http.StatusUnauthorized)
				return
			}
		default:
			c.AbortWithStatus(http.StatusUnauthorized)
		}

		c.Next()
	}
}
