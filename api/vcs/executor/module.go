package executor

import (
	"getsturdy.com/api/pkg/di"
	"getsturdy.com/api/pkg/logger"
	"getsturdy.com/api/vcs/provider"
)

func Module(c *di.Container) {
	c.Import(logger.Module)
	c.Import(provider.Module)
	c.Register(NewProvider)
}
