//go:build !cloud
// +build !cloud

package module

import (
	"mash/pkg/analytics/proxy"
	"mash/pkg/di"
)

func Module(c *di.Container) {
	c.Import(proxy.Module)
}
