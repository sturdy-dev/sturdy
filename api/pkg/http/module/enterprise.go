//go:build enterprise
// +build enterprise

package http

import (
	"net/http"

	"getsturdy.com/api/pkg/di"
	httpx "getsturdy.com/api/pkg/http"
	"getsturdy.com/api/pkg/http/enterprise/selfhosted"

	"github.com/gin-gonic/gin"
)

func Module(c *di.Container) {
	c.Register(httpx.ProvideHandler)
	c.Register(selfhosted.ProvideHandler)
	c.Register(func(e *selfhosted.Engine) http.Handler {
		return (*gin.Engine)(e)
	})
	c.Register(httpx.ProvideServer)
}
