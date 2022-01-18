package module

import (
	"mash/pkg/auth/service"
	"mash/pkg/di"
)

func Module(c *di.Container) {
	c.Import(service.Module)
}
