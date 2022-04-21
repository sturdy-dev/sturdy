package graphql

import (
	services_auth "getsturdy.com/api/pkg/auth/service"
	"getsturdy.com/api/pkg/di"
	service_land "getsturdy.com/api/pkg/land/enterprise/service"
	service_users "getsturdy.com/api/pkg/users/service/module"
	graphql_workspaces "getsturdy.com/api/pkg/workspaces/graphql"
	service_workspaces "getsturdy.com/api/pkg/workspaces/service"
)

func Module(c *di.Container) {
	c.Import(services_auth.Module)
	c.Import(service_land.Module)
	c.Import(service_users.Module)
	c.Import(service_workspaces.Module)
	c.Import(graphql_workspaces.Module)
	c.Register(NewResolver)
}
