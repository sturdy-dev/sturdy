//go:build cloud
// +build cloud

package http

import (
	"getsturdy.com/api/pkg/di"
	"getsturdy.com/api/pkg/http/cloud"
	"getsturdy.com/api/pkg/http/enterprise"
	"getsturdy.com/api/pkg/http/oss"
	"net/http"
)

func Module(c *di.Container) {
	c.Register(oss.ProvideHandler)
	c.Register(enterprise.ProvideHandler)
	c.Register(cloud.ProvideHandler, new(http.Handler))
	c.Register(ProvideServer)
}
