package workers

import (
	"getsturdy.com/api/pkg/di"
)

func Module(c *di.Container) {
	c.Register(NewImporterQueue)
	c.Register(NewClonerQueue)
	c.Register(NewWebhooksQueue)
}
