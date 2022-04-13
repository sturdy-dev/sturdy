package service

import (
	"getsturdy.com/api/pkg/di"
	"getsturdy.com/api/pkg/events/v2"
	"getsturdy.com/api/pkg/logger"
	db_statuses "getsturdy.com/api/pkg/statuses/db"
)

func Module(c *di.Container) {
	c.Import(logger.Module)
	c.Import(db_statuses.Module)
	c.Import(events.Module)
	c.Register(New)
}
