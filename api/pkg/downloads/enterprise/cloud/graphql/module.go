package graphql

import (
	service_auth "getsturdy.com/api/pkg/auth/service"
	"getsturdy.com/api/pkg/di"
	service_downloads "getsturdy.com/api/pkg/downloads/enterprise/cloud/service"
	service_snapshots "getsturdy.com/api/pkg/snapshots/service"
)

func Module(c *di.Container) {
	c.Import(service_downloads.Module)
	c.Import(service_auth.Module)
	c.Import(service_snapshots.Module)
	c.Register(New)
}
