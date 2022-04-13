package cloud

import (
	"getsturdy.com/api/pkg/aws"
	configuration "getsturdy.com/api/pkg/configuration/module"
	"getsturdy.com/api/pkg/di"
	"getsturdy.com/api/pkg/logger"
)

func Module(c *di.Container) {
	c.Import(configuration.Module)
	c.Import(aws.Module)
	c.Import(logger.Module)
	c.Register(New)
}
