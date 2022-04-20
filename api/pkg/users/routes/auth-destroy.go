package routes

import (
	"net/http"

	"getsturdy.com/api/pkg/auth"

	"github.com/gin-gonic/gin"
)

func AuthDestroy(c *gin.Context) {
	auth.RemoveAuthCookie(c.Writer)
	c.Status(http.StatusOK)
}
