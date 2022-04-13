package service

import (
	db_ci "getsturdy.com/api/pkg/ci/db"
	configuration "getsturdy.com/api/pkg/configuration/module"
	"getsturdy.com/api/pkg/di"
	db_integrations "getsturdy.com/api/pkg/integrations/db"
	service_jwt "getsturdy.com/api/pkg/jwt/service"
	"getsturdy.com/api/pkg/logger"
	"getsturdy.com/api/pkg/snapshots/snapshotter"
	service_statuses "getsturdy.com/api/pkg/statuses/service"
	"getsturdy.com/api/vcs/executor"
)

func Module(c *di.Container) {
	c.Import(logger.Module)
	c.Import(executor.Module)
	c.Import(db_integrations.Module)
	c.Import(db_ci.Module)
	c.Import(configuration.Module)
	c.Import(service_statuses.Module)
	c.Import(service_jwt.Module)
	c.Import(snapshotter.Module)
	c.Register(New)
}
