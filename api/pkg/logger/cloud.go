//go:build cloud
// +build cloud

package logger

import (
	configuration "getsturdy.com/api/pkg/configuration/module"
	"getsturdy.com/api/pkg/di"
	"getsturdy.com/api/pkg/logger/enterprise/cloud"
)

func Module(c *di.Container) {
	c.Import(configuration.Module)
	c.Import(cloud.Module)
	c.Register(New)
}
