//go:build enterprise || cloud
// +build enterprise cloud

package http

import (
	"mash/pkg/di"
	"mash/pkg/http/enterprise"
	"mash/pkg/http/oss"
)

func Module(c *di.Container) {
	c.Register(oss.ProvideHandler)
	c.Register(enterprise.ProvideHandler)
	c.Register(ProvideServer)
}
