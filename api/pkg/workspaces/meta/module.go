package meta

import (
	"getsturdy.com/api/pkg/di"
	"getsturdy.com/api/pkg/events"
	"getsturdy.com/api/pkg/logger"
	db_workspaces "getsturdy.com/api/pkg/workspaces/db"
)

func Module(c *di.Container) {
	c.Import(logger.Module)
	c.Import(events.Module)
	c.Import(db_workspaces.Module)
	c.Register(NewWriterWithEvents)
}
