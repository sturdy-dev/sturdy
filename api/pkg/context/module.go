package context

import (
	"context"

	"getsturdy.com/api/pkg/di"
)

func Module(c *di.Container) {
	c.Register(func() context.Context {
		return context.Background()
	})
}
