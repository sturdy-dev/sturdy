package proxy

import (
	"getsturdy.com/api/pkg/di"
	global_installations "getsturdy.com/api/pkg/installations/global"
)

func Module(c *di.Container) {
	c.Import(global_installations.Module)
	c.Register(NewClient)
}
