package module

import (
	"getsturdy.com/api/pkg/di"
	"getsturdy.com/api/vcs/executor"
	"getsturdy.com/api/vcs/provider"
)

func Module(c *di.Container) {
	c.Import(executor.Module)
	c.Import(provider.Module)
}
