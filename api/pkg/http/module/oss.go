//go:build !enterprise && !cloud
// +build !enterprise,!cloud

package http

import (
	"net/http"

	"getsturdy.com/api/pkg/di"
	httpx "getsturdy.com/api/pkg/http"

	"github.com/gin-gonic/gin"
)

func Module(c *di.Container) {
	c.Register(httpx.ProvideHandler)
	c.Register(func(e *httpx.Engine) http.Handler {
		return (*gin.Engine)(e)
	})
	c.Register(httpx.ProvideServer)
}
