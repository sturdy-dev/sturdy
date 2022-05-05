package module

import (
	"testing"

	"getsturdy.com/api/pkg/di"
	"getsturdy.com/api/pkg/logger"
	"getsturdy.com/api/pkg/queue"
)

func TestModule(t *testing.T) di.Module {
	return func(c *di.Container) {
		c.Register(func() queue.Queue { return queue.NewInMemory(logger.NewTest(t)).Sync() })
	}
}
