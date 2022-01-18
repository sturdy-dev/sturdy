//go:build enterprise
// +build enterprise

package http

import (
	"net/http"

	"mash/pkg/di"
	"mash/pkg/http/enterprise"
	"mash/pkg/http/oss"

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
