//go:build cloud
// +build cloud

package http

import (
	"net/http"

	"getsturdy.com/api/pkg/di"
	httpx "getsturdy.com/api/pkg/http"
	"getsturdy.com/api/pkg/http/cloud"
	"getsturdy.com/api/pkg/http/enterprise"
)

func Module(c *di.Container) {
	c.Register(httpx.ProvideHandler)
	c.Register(enterprise.ProvideHandler)
	c.Register(cloud.ProvideHandler, new(http.Handler))
	c.Register(httpx.ProvideServer)
}
