//go:build !enterprise
// +build !enterprise

package http

import (
	"net/http"

	"mash/pkg/di"
	"mash/pkg/http/oss"
)

func Module(c *di.Container) {
	c.Register(oss.ProvideHandler, new(http.Handler))
	c.Register(ProvideServer)
}
