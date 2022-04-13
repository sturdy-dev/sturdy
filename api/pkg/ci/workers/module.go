package workers

import (
	service_ci "getsturdy.com/api/pkg/ci/service"
	"getsturdy.com/api/pkg/di"
	"getsturdy.com/api/pkg/logger"
	queue "getsturdy.com/api/pkg/queue/module"
)

func Module(c *di.Container) {
	c.Import(logger.Module)
	c.Import(queue.Module)
	c.Import(service_ci.Module)
	c.Register(New)
}
