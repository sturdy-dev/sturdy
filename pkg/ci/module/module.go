package module

import (
	"mash/pkg/ci/db"
	"mash/pkg/ci/service"
	"mash/pkg/ci/workers"
	"mash/pkg/di"
)

func Module(c *di.Container) {
	c.Import(db.Module)
	c.Import(service.Module)
	c.Import(workers.Module)
}
