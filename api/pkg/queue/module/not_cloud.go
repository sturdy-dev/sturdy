//go:build !cloud
// +build !cloud

package queue

import (
	"getsturdy.com/api/pkg/di"
	"getsturdy.com/api/pkg/queue"
)

func Module(c *di.Container) {
	c.Register(queue.NewInMemory, new(queue.Queue))
}
