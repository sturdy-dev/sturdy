package graphql

import (
	service_auth "getsturdy.com/api/pkg/auth/service"
	"getsturdy.com/api/pkg/di"
	service_downloads "getsturdy.com/api/pkg/downloads/enterprise/cloud/service"
	"getsturdy.com/api/pkg/snapshots/snapshotter"
)

func Module(c *di.Container) {
	c.Import(service_downloads.Module)
	c.Import(service_auth.Module)
	c.Import(snapshotter.Module)
	c.Register(New)
}
