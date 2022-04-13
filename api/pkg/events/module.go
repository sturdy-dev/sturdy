package events

import (
	db_codebases "getsturdy.com/api/pkg/codebases/db"
	"getsturdy.com/api/pkg/di"
	"getsturdy.com/api/pkg/logger"
	db_organizations "getsturdy.com/api/pkg/organization/db"
	db_workspaces "getsturdy.com/api/pkg/workspaces/db"
)

func Module(c *di.Container) {
	c.Import(logger.Module)
	c.Import(db_codebases.Module)
	c.Import(db_organizations.Module)
	c.Import(db_workspaces.Module)
	c.Register(NewInMemory)
	c.Register(NewSender)
	c.Register(func(e EventReadWriter) EventReader {
		return e
	})
}
