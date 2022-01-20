//go:build cloud
// +build cloud

package http

import (
	"mash/pkg/di"
	"mash/pkg/http/cloud"
	"mash/pkg/http/enterprise"
	"mash/pkg/http/oss"
	"net/http"
)

func Module(c *di.Container) {
	c.Register(oss.ProvideHandler)
	c.Register(enterprise.ProvideHandler)
	c.Register(cloud.ProvideHandler, new(http.Handler))
	c.Register(ProvideServer)
}
