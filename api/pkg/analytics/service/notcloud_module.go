//go:build !cloud
// +build !cloud

package service

import (
	"getsturdy.com/api/pkg/analytics/proxy"
	"getsturdy.com/api/pkg/di"
	"getsturdy.com/api/pkg/logger"
)

func Module(c *di.Container) {
	c.Import(logger.Module)
	c.Import(proxy.Module)
	c.Register(New)
}
