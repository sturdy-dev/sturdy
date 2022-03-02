package db

import "getsturdy.com/api/pkg/di"

func Module(c *di.Container) {
	c.Register(NewRepo)
}

func TestModule(c *di.Container) {
	c.Register(NewInMemoryRepo)
}
