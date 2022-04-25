package datamigrations

import (
	service_changes "getsturdy.com/api/pkg/changes/service"
	service_codebases "getsturdy.com/api/pkg/codebases/service"
	"getsturdy.com/api/pkg/di"
	"getsturdy.com/api/pkg/logger"
	"getsturdy.com/api/pkg/review/db"
)

func Module(c *di.Container) {
	c.Import(db.Module)
	c.Import(service_changes.Module)
	c.Import(service_codebases.Module)
	c.Import(logger.Module)
	c.Register(NewService)
}
