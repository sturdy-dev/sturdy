package db

import (
	"getsturdy.com/api/pkg/di"
)

func Module(c *di.Container) {
	c.Register(NewRepo)
	c.Register(func(repo Repository) WorkspaceReader {
		return repo
	})
}
