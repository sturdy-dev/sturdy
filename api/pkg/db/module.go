package db

import (
	module_configuration "getsturdy.com/api/pkg/configuration/module"
	"getsturdy.com/api/pkg/di"
)

func Module(c *di.Container) {
	c.Import(module_configuration.Module)
	c.Register(FromConfiguration)
}
