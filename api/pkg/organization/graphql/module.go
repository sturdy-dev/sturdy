package graphql

import (
	service_auth "getsturdy.com/api/pkg/auth/service"
	graphql_author "getsturdy.com/api/pkg/author/graphql"
	graphql_codebases "getsturdy.com/api/pkg/codebases/graphql"
	service_codebase "getsturdy.com/api/pkg/codebases/service"
	"getsturdy.com/api/pkg/di"
	"getsturdy.com/api/pkg/events/v2"
	"getsturdy.com/api/pkg/graphql/resolvers"
	graphql_licenses "getsturdy.com/api/pkg/licenses/graphql"
	"getsturdy.com/api/pkg/logger"
	service_organization "getsturdy.com/api/pkg/organization/service"
	service_user "getsturdy.com/api/pkg/users/service"
)

func Module(c *di.Container) {
	c.Import(service_organization.Module)
	c.Import(service_auth.Module)
	c.Import(service_user.Module)
	c.Import(service_codebase.Module)
	c.Import(graphql_author.Module)
	c.Import(graphql_licenses.Module)
	c.Import(graphql_codebases.Module)
	c.Import(logger.Module)
	c.Import(events.Module)
	c.Register(New)

	// populate cyclic resolver
	c.Import(resolvers.Module)
	c.Decorate(func(rp *resolvers.OrganizationRootResolver, rv resolvers.OrganizationRootResolver) *resolvers.OrganizationRootResolver {
		*rp = rv
		return &rv
	})
}
