//go:build !enterprise && !cloud
// +build !enterprise,!cloud

package module

import (
	"net/http"

	"getsturdy.com/api/pkg/di"
	"getsturdy.com/api/pkg/http/handler"
	"github.com/gin-gonic/gin"
)

func Module(c *di.Container) {
	c.Import(handler.Module)
	c.Register(func(e *handler.Engine) http.Handler {
		return (*gin.Engine)(e)
	})
}
