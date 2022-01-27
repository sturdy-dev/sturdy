//go:build cloud
// +build cloud

package http

import (
	"net/http"

	"getsturdy.com/api/pkg/di"
	httpx "getsturdy.com/api/pkg/http"
	"getsturdy.com/api/pkg/http/enterprise/cloud"
	"getsturdy.com/api/pkg/http/enterprise/selfhosted"
)

func Module(c *di.Container) {
	c.Register(httpx.ProvideHandler)
	c.Register(selfhosted.ProvideHandler)
	c.Register(cloud.ProvideHandler, new(http.Handler))
	c.Register(httpx.ProvideServer)
}
