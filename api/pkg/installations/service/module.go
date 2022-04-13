package service

import (
	"getsturdy.com/api/pkg/di"
	"getsturdy.com/api/pkg/installations/db"
)

func Module(c *di.Container) {
	c.Import(db.Module)
	c.Register(New)
}
