package module

import (
	"getsturdy.com/api/pkg/auth/service"
	"getsturdy.com/api/pkg/di"
)

func Module(c *di.Container) {
	c.Import(service.Module)
}
