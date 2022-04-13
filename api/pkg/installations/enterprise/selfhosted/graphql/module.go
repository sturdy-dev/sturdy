package graphql

import (
	"getsturdy.com/api/pkg/di"
	"getsturdy.com/api/pkg/graphql/resolvers"
	"getsturdy.com/api/pkg/installations/enterprise/selfhosted/service"
	"getsturdy.com/api/pkg/installations/graphql"
	graphql_licenses "getsturdy.com/api/pkg/licenses/graphql"
)

func Module(c *di.Container) {
	c.Import(graphql.Module)
	c.Import(service.Module)
	c.Import(graphql_licenses.Module)
	c.Register(New, new(resolvers.InstallationsRootResolver))
}
