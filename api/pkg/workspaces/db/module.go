package db

import (
	"getsturdy.com/api/pkg/db"
	"getsturdy.com/api/pkg/di"
)

func Module(c *di.Container) {
	c.Import(db.Module)
	c.Register(NewRepo)
	c.Register(func(repo Repository) WorkspaceReader {
		return repo
	})
}

func TestModule(c *di.Container) {
	c.Register(NewMemory)
	c.Register(func(repo Repository) WorkspaceReader {
		return repo
	})
	c.Register(func(repo Repository) WorkspaceWriter {
		return repo
	})

}
