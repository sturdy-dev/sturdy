//go:build !enterprise && !cloud
// +build !enterprise,!cloud

package http

import (
	"net/http"

	"mash/pkg/di"
	"mash/pkg/http/oss"

	"github.com/gin-gonic/gin"
)

func Module(c *di.Container) {
	c.Register(oss.ProvideHandler)
	c.Register(func(e *oss.Engine) http.Handler {
		return (*gin.Engine)(e)
	})
	c.Register(ProvideServer)
}
