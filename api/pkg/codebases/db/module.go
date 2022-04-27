package db

import (
	"getsturdy.com/api/pkg/db"
	"getsturdy.com/api/pkg/di"
)

func Module(c *di.Container) {
	c.Import(db.Module)
	c.Register(NewCodebaseUserRepo)
	c.Register(NewRepo)
}

func TestModule(c *di.Container) {
	c.Register(NewInMemoryCodebaseUserRepo)
	c.Register(NewMemory)
}
