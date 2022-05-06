package meta

import (
	"getsturdy.com/api/pkg/di"
	events "getsturdy.com/api/pkg/events/v2"
	worker_snapshotter "getsturdy.com/api/pkg/snapshots/worker"
	db_workspaces "getsturdy.com/api/pkg/workspaces/db"
)

func Module(c *di.Container) {
	c.Import(db_workspaces.Module)
	c.Import(events.Module)
	c.Import(worker_snapshotter.Module)
	c.Register(NewViewUpdatedFunc)
}
