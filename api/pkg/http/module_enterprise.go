//go:build enterprise
// +build enterprise

package http

import (
	"net/http"

	"getsturdy.com/api/pkg/di"
	"getsturdy.com/api/pkg/http/enterprise"
	"getsturdy.com/api/pkg/http/oss"

	"github.com/gin-gonic/gin"
)

func Module(c *di.Container) {
	c.Register(oss.ProvideHandler)
	c.Register(enterprise.ProvideHandler)
	c.Register(func(e *enterprise.Engine) http.Handler {
		return (*gin.Engine)(e)
	})
	c.Register(ProvideServer)
}
