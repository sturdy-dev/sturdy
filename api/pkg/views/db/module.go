package db

import (
	"getsturdy.com/api/pkg/di"
	"getsturdy.com/api/pkg/users/db"
)

func Module(c *di.Container) {
	c.Import(db.Module)
	c.Register(NewRepo)
}

func TestModule(c *di.Container) {
	c.Register(NewInMemoryViewRepo)
}
