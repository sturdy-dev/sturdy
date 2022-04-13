package service

import (
	"getsturdy.com/api/pkg/di"
	"getsturdy.com/api/pkg/events"
	db_watchers "getsturdy.com/api/pkg/workspaces/watchers/db"
)

func Module(c *di.Container) {
	c.Import(db_watchers.Module)
	c.Import(events.Module)
	c.Register(New)
}
