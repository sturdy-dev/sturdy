package cloud

import (
	"getsturdy.com/api/pkg/aws"
	"getsturdy.com/api/pkg/di"
	"getsturdy.com/api/pkg/logger"
)

func Module(c *di.Container) {
	c.Import(logger.Module)
	c.Import(aws.Module)
	c.Register(New)
}
