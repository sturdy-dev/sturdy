package graphql

import (
	service_auth "getsturdy.com/api/pkg/auth/service"
	graphql_buildkite "getsturdy.com/api/pkg/buildkite/graphql/module"
	service_change "getsturdy.com/api/pkg/changes/service"
	service_ci "getsturdy.com/api/pkg/ci/service"
	"getsturdy.com/api/pkg/di"
	"getsturdy.com/api/pkg/graphql/resolvers"
	graphql_statuses "getsturdy.com/api/pkg/statuses/graphql/module"
	service_workspaces "getsturdy.com/api/pkg/workspaces/service"
)

func Module(c *di.Container) {
	c.Import(service_ci.Module)
	c.Import(service_auth.Module)
	c.Import(service_change.Module)
	c.Import(service_workspaces.Module)
	c.Import(graphql_statuses.Module)
	c.Import(graphql_buildkite.Module)
	c.Register(NewRootResolver)

	// populate cyclic resolver
	c.Decorate(func(rv *resolvers.IntegrationRootResolver, rp resolvers.IntegrationRootResolver) *resolvers.IntegrationRootResolver {
		*rv = rp
		return &rp
	})
}
