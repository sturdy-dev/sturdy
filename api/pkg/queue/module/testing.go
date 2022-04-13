package module

import (
	"getsturdy.com/api/pkg/di"
	"getsturdy.com/api/pkg/queue"
)

func TestModule(c *di.Container) {
	c.Register(queue.NewSync, new(queue.Queue))
}
