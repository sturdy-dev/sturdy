package events

import "getsturdy.com/api/pkg/di"

func Module(c *di.Container) {
	c.Register(NewInMemory)
	c.Register(NewSender)
	c.Register(func(e EventReadWriter) EventReader {
		return e
	})
}
