package proxy

import (
	configuration "getsturdy.com/api/pkg/configuration/module"
	"getsturdy.com/api/pkg/di"
	global_installations "getsturdy.com/api/pkg/installations/global"
	"getsturdy.com/api/pkg/logger"
)

func Module(c *di.Container) {
	c.Import(configuration.Module)
	c.Import(global_installations.Module)
	c.Import(logger.Module)
	c.Register(NewClient)
}
