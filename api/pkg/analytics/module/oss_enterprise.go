//go:build !cloud
// +build !cloud

package module

import (
	"getsturdy.com/api/pkg/analytics/proxy"
	"getsturdy.com/api/pkg/analytics/service"
	"getsturdy.com/api/pkg/di"
)

func Module(c *di.Container) {
	c.Import(proxy.Module)
	c.Import(service.Module)
}
