package inmemory

import (
	"getsturdy.com/api/pkg/di"
)

func TestModule(c *di.Container) {
	c.Register(NewInMemoryAclRepo)
}
