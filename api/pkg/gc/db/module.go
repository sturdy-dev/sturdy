package db

import (
	"getsturdy.com/api/pkg/di"
	"getsturdy.com/api/pkg/notification/db"
)

func Module(c *di.Container) {
	c.Import(db.Module)
	c.Register(NewRepository)
}
