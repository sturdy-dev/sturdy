package service

import (
	"getsturdy.com/api/pkg/di"
	"getsturdy.com/api/pkg/events/v2"
	"getsturdy.com/api/pkg/logger"
	db_watchers "getsturdy.com/api/pkg/workspaces/watchers/db"
)

func Module(c *di.Container) {
	c.Import(db_watchers.Module)
	c.Import(events.Module)
	c.Import(logger.Module)
	c.Register(New)
}
