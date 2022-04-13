package events

import (
	db_codebases "getsturdy.com/api/pkg/codebases/db"
	"getsturdy.com/api/pkg/di"
	"getsturdy.com/api/pkg/logger"
	db_organization "getsturdy.com/api/pkg/organization/db"
	db_workspaces "getsturdy.com/api/pkg/workspaces/db"
)

func Module(c *di.Container) {
	c.Import(logger.Module)
	c.Import(db_codebases.Module)
	c.Import(db_organization.Module)
	c.Import(db_workspaces.Module)

	c.Register(New)
	c.Register(NewPublisher)
	c.Register(NewSubscriber)
}
