package service

import (
	service_buildkite "getsturdy.com/api/pkg/buildkite/service/module"
	db_ci "getsturdy.com/api/pkg/ci/db"
	configuration "getsturdy.com/api/pkg/configuration/module"
	"getsturdy.com/api/pkg/di"
	service_github "getsturdy.com/api/pkg/github/service/module"
	db_integrations "getsturdy.com/api/pkg/integrations/db"
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
	c.Import(service_buildkite.Module)
	c.Import(service_github.Module)
	c.Register(New)
}
