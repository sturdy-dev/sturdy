package module

import (
	"getsturdy.com/api/pkg/di"
	"getsturdy.com/api/pkg/sync/service"
)

func Module(c *di.Container) {
	c.Import(service.Module)
}
