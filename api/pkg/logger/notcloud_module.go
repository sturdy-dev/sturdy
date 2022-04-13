//go:build !cloud
// +build !cloud

package logger

import (
	configuration "getsturdy.com/api/pkg/configuration/module"
	"getsturdy.com/api/pkg/di"
	"getsturdy.com/api/pkg/logger/proxy"
)

func Module(c *di.Container) {
	c.Import(proxy.Module)
	c.Import(configuration.Module)
	c.Register(New)
}
