package worker

import "getsturdy.com/api/pkg/di"

func Module(c *di.Container) {
	c.Register(New)
}

func TestingModule(c *di.Container) {
	c.Register(NewSync)
}
