package graphql

import (
	"getsturdy.com/api/pkg/di"
	service_licenses "getsturdy.com/api/pkg/licenses/enterprise/cloud/service"
)

func Module(c *di.Container) {
	c.Import(service_licenses.Module)
	c.Register(New)
}
