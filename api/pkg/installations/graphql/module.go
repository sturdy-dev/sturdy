package graphql

import (
	"getsturdy.com/api/pkg/di"
	"getsturdy.com/api/pkg/installations/service"
	graphql_licences "getsturdy.com/api/pkg/licenses/graphql"
	service_organizations "getsturdy.com/api/pkg/organization/service"
	service_users "getsturdy.com/api/pkg/users/service/module"
)

func Module(c *di.Container) {
	c.Import(service.Module)
	c.Import(graphql_licences.Module)
	c.Import(service_organizations.Module)
	c.Import(service_users.Module)
	c.Register(New)
}
