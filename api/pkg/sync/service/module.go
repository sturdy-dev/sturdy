package service

import (
	"getsturdy.com/api/pkg/di"
	"getsturdy.com/api/pkg/events"
	"getsturdy.com/api/pkg/logger"
	"getsturdy.com/api/pkg/snapshots/snapshotter"
	db_view "getsturdy.com/api/pkg/view/db"
	db_workspaces "getsturdy.com/api/pkg/workspaces/db"
	"getsturdy.com/api/vcs/executor"
)

func Module(c *di.Container) {
	c.Import(logger.Module)
	c.Import(executor.Module)
	c.Import(db_view.Module)
	c.Import(db_workspaces.Module)
	c.Import(snapshotter.Module)
	c.Import(events.Module)
	c.Register(New)
}
