package gitserver

import (
	service_codebase "getsturdy.com/api/pkg/codebases/service"
	configuration "getsturdy.com/api/pkg/configuration/module"
	"getsturdy.com/api/pkg/di"
	service_jwt "getsturdy.com/api/pkg/jwt/service"
	"getsturdy.com/api/pkg/logger"
	service_servicetokens "getsturdy.com/api/pkg/servicetokens/service"
	"getsturdy.com/api/vcs/executor"
)

func Module(c *di.Container) {
	c.Import(logger.Module)
	c.Import(configuration.Module)
	c.Import(service_servicetokens.Module)
	c.Import(service_jwt.Module)
	c.Import(service_codebase.Module)
	c.Import(executor.Module)
	c.Register(New)
}
