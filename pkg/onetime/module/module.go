package module

import (
	"mash/pkg/di"
	"mash/pkg/onetime/db"
	"mash/pkg/onetime/service"
)

func Module(c *di.Container) {
	c.Import(db.Module)
	c.Import(service.Module)
}
