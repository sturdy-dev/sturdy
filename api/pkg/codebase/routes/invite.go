package routes

import (
	"database/sql"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	service_auth "getsturdy.com/api/pkg/auth/service"
	service_codebase "getsturdy.com/api/pkg/codebase/service"
)

type InviteUserRequest struct {
	UserEmail string `json:"user_email" binding:"required"`
}

func Invite(
	codebaseService *service_codebase.Service,
	authService *service_auth.Service,
) func(c *gin.Context) {
	return func(c *gin.Context) {
		codebaseID := c.Param("id")

		cb, err := codebaseService.GetByID(c.Request.Context(), codebaseID)
		if err != nil && errors.Is(err, sql.ErrNoRows) {
			c.AbortWithStatus(http.StatusNotFound)
			return
		} else if err != nil {
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		ctx := c.Request.Context()

		if err := authService.CanWrite(ctx, cb); err != nil {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		var request InviteUserRequest
		if err := c.BindJSON(&request); err != nil {
			c.AbortWithStatus(http.StatusBadRequest)
			return
		}

		member, err := codebaseService.AddUserByEmail(ctx, cb.ID, request.UserEmail)
		if err != nil {
			c.AbortWithStatus(http.StatusBadRequest)
			return
		}
		c.JSON(http.StatusOK, member)
	}
}
