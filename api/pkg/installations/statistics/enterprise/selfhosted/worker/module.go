package worker

import (
	"getsturdy.com/api/pkg/di"
	"getsturdy.com/api/pkg/installations/statistics/enterprise/selfhosted/publisher"
	"getsturdy.com/api/pkg/installations/statistics/enterprise/selfhosted/service"
	"getsturdy.com/api/pkg/logger"
)

func Module(c *di.Container) {
	c.Import(logger.Module)
	c.Import(service.Module)
	c.Import(publisher.Module)
	c.Register(New)
}
