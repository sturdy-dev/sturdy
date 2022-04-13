package service

import (
	"getsturdy.com/api/pkg/activity/db"
	"getsturdy.com/api/pkg/di"
	"getsturdy.com/api/pkg/events"
)

func Module(c *di.Container) {
	c.Import(db.Module)
	c.Import(events.Module)
	c.Register(New)
}
