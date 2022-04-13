package service

import (
	"getsturdy.com/api/pkg/di"
	"getsturdy.com/api/pkg/events"
	db_presence "getsturdy.com/api/pkg/presence/db"
)

func Module(c *di.Container) {
	c.Import(db_presence.Module)
	c.Import(events.Module)
	c.Register(New)
}
