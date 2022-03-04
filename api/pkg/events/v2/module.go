package events

import (
	"getsturdy.com/api/pkg/di"
)

func Module(c *di.Container) {
	c.Register(New)
	c.Register(NewPublisher)
	c.Register(NewSubscriber)
}
