//go:build enterprise
// +build enterprise

package http

import (
	"net/http"

	"getsturdy.com/api/pkg/di"
	"getsturdy.com/api/pkg/http/handler/enterprise/selfhosted"

	"github.com/gin-gonic/gin"
)

func Module(c *di.Container) {
	c.Import(selfhosted.Module)
	c.Register(func(e *selfhosted.Engine) http.Handler {
		return (*gin.Engine)(e)
	})
}
