package worker

import (
	"getsturdy.com/api/pkg/di"
	"getsturdy.com/api/pkg/logger"
	queue "getsturdy.com/api/pkg/queue/module"
	"getsturdy.com/api/pkg/snapshots/snapshotter"
)

func Module(c *di.Container) {
	c.Import(logger.Module)
	c.Import(queue.Module)
	c.Import(snapshotter.Module)
	c.Register(New)
}

func TestModule(c *di.Container) {
	c.Register(NewSync)
}
