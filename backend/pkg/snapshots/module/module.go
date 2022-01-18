package module

import (
	"mash/pkg/di"
	"mash/pkg/snapshots/db"
	"mash/pkg/snapshots/snapshotter"
	"mash/pkg/snapshots/worker"
)

func Module(c *di.Container) {
	c.Import(db.Module)
	c.Import(snapshotter.Module)
	c.Import(worker.Module)
}
