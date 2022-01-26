package module

import (
	"getsturdy.com/api/pkg/di"
	"getsturdy.com/api/pkg/user/service"
)

func Module(c *di.Container) {
	c.Register(service.New)
}
