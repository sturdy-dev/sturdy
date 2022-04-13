package worker

import (
	"getsturdy.com/api/pkg/di"
	gc_service "getsturdy.com/api/pkg/gc/service"
	"getsturdy.com/api/pkg/logger"
	queue "getsturdy.com/api/pkg/queue/module"
)

func Module(c *di.Container) {
	c.Import(logger.Module)
	c.Import(queue.Module)
	c.Import(gc_service.Module)
	c.Register(New)
}
