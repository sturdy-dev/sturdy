package service

import (
	"getsturdy.com/api/pkg/di"
	db_gc "getsturdy.com/api/pkg/gc/db"
	"getsturdy.com/api/pkg/logger"
	db_snapshots "getsturdy.com/api/pkg/snapshots/db"
	service_suggestions "getsturdy.com/api/pkg/suggestions/service"
	db_view "getsturdy.com/api/pkg/view/db"
	db_workspaces "getsturdy.com/api/pkg/workspaces/db"
	"getsturdy.com/api/vcs/executor"
)

func Module(c *di.Container) {
	c.Import(logger.Module)
	c.Import(db_gc.Module)
	c.Import(db_view.Module)
	c.Import(db_snapshots.Module)
	c.Import(db_workspaces.Module)
	c.Import(service_suggestions.Module)
	c.Import(executor.Module)
	c.Register(New)
}
