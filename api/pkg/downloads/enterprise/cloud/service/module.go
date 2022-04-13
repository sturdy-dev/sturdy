package service

import (
	"getsturdy.com/api/pkg/aws"
	service_change "getsturdy.com/api/pkg/changes/service"
	configuration "getsturdy.com/api/pkg/configuration/module"
	"getsturdy.com/api/pkg/di"
	"getsturdy.com/api/pkg/logger"
	"getsturdy.com/api/vcs/executor"
)

func Module(c *di.Container) {
	c.Import(logger.Module)
	c.Import(configuration.Module)
	c.Import(executor.Module)
	c.Import(service_change.Module)
	c.Import(aws.Module)
	c.Register(New)
}
