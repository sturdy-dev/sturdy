package routes

import (
	"getsturdy.com/api/pkg/di"
	"getsturdy.com/api/pkg/logger"
	"getsturdy.com/api/pkg/remote/enterprise/service"
)

func Module(c *di.Container) {
	c.Import(service.Module)
	c.Import(logger.Module)
	c.Register(TriggerSyncCodebaseWebhook)
}
