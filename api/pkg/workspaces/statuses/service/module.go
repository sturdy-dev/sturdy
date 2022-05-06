package service

import (
	"getsturdy.com/api/pkg/di"
	service_snapshots "getsturdy.com/api/pkg/snapshots/service"
	service_statuses "getsturdy.com/api/pkg/statuses/service"
)

func Module(c *di.Container) {
	c.Import(service_snapshots.Module)
	c.Import(service_statuses.Module)
	c.Register(New)
}
