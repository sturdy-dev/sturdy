package module

import (
	"getsturdy.com/api/pkg/di"
	"getsturdy.com/api/pkg/queue"

	"go.uber.org/zap"
)

func TestModule(c *di.Container) {
	c.Register(func() queue.Queue { return queue.NewInMemory(zap.NewNop()).Sync() })
}
