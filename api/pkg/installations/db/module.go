package db

import (
	"getsturdy.com/api/pkg/db"
	"getsturdy.com/api/pkg/di"
)

func Module(c *di.Container) {
	c.Import(db.Module)
	c.Register(NewDatabase)
}
