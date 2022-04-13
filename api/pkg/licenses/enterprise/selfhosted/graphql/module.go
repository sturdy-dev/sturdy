package graphql

import (
	"getsturdy.com/api/pkg/di"
	service_installations "getsturdy.com/api/pkg/installations/service"
	"getsturdy.com/api/pkg/licenses/enterprise/selfhosted/validator"
)

func Module(c *di.Container) {
	c.Import(service_installations.Module)
	c.Import(validator.Module)
	c.Register(New)
}
