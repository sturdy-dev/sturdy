//go:build !cloud
// +build !cloud

package module

import (
	"getsturdy.com/api/pkg/di"
	"getsturdy.com/api/pkg/logger"
	"getsturdy.com/api/pkg/logger/proxy"
)

func Module(c *di.Container) {
	c.Import(proxy.Module)
	c.Register(logger.New)
}
