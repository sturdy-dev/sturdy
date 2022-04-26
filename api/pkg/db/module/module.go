package module

import (
	"getsturdy.com/api/pkg/configuration"
	module_configuration "getsturdy.com/api/pkg/configuration/module"
	"getsturdy.com/api/pkg/db"
	"getsturdy.com/api/pkg/di"
)

func Module(c *di.Container) {
	c.Import(module_configuration.Module)
	c.Register(db.FromConfiguration)
}

func TestModule(c *di.Container) {
	c.Import(configuration.TestModule)
}
