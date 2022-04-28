package graphql

import (
	"getsturdy.com/api/pkg/di"
	service_snapshots "getsturdy.com/api/pkg/snapshots/service"
)

func Module(c *di.Container) {
	c.Import(service_snapshots.Module)
	c.Register(NewRoot)
}
