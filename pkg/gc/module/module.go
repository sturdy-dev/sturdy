package module

import (
	"mash/pkg/di"
	"mash/pkg/gc/db"
	"mash/pkg/gc/worker"
)

func Module(c *di.Container) {
	c.Import(db.Module)
	c.Import(worker.Module)
}
