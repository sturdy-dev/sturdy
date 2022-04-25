package service

import (
	db_ci "getsturdy.com/api/pkg/ci/db"
	configuration "getsturdy.com/api/pkg/configuration/module"
	"getsturdy.com/api/pkg/di"
	db_integrations "getsturdy.com/api/pkg/integrations/db"
	integration_providers "getsturdy.com/api/pkg/integrations/providers/module"
	service_jwt "getsturdy.com/api/pkg/jwt/service"
	"getsturdy.com/api/pkg/logger"
	service_snapshots "getsturdy.com/api/pkg/snapshots/service"
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
	c.Import(service_snapshots.Module)
	c.Import(integration_providers.Module)
	c.Register(New)
}
