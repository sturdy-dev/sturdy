package global

import (
	"getsturdy.com/api/pkg/context"
	"getsturdy.com/api/pkg/di"
	"getsturdy.com/api/pkg/installations/service"
)

func Module(c *di.Container) {
	c.Import(context.Module)
	c.Import(service.Module)
	c.Register(New)
}
