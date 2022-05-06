package graphql

import (
	service_auth "getsturdy.com/api/pkg/auth/service"
	graphql_changes "getsturdy.com/api/pkg/changes/graphql"
	service_changes "getsturdy.com/api/pkg/changes/service"
	"getsturdy.com/api/pkg/di"
	"getsturdy.com/api/pkg/events/v2"
	graphql_github_pr "getsturdy.com/api/pkg/github/graphql/pr"
	"getsturdy.com/api/pkg/graphql/resolvers"
	"getsturdy.com/api/pkg/logger"
	service_snapshots "getsturdy.com/api/pkg/snapshots/service"
	service_statuses "getsturdy.com/api/pkg/statuses/service"
	service_workspace "getsturdy.com/api/pkg/workspaces/service"
	service_workspace_statuses "getsturdy.com/api/pkg/workspaces/statuses/service"
)

func Module(c *di.Container) {
	c.Import(logger.Module)
	c.Import(service_statuses.Module)
	c.Import(service_changes.Module)
	c.Import(service_workspace.Module)
	c.Import(service_auth.Module)
	c.Import(graphql_changes.Module)
	c.Import(graphql_github_pr.Module)
	c.Import(resolvers.Module)
	c.Import(events.Module)
	c.Import(service_snapshots.Module)
	c.Import(service_workspace_statuses.Module)
	c.Register(New, new(resolvers.StatusesRootResolver))

	// populate cyclic resolver
	c.Decorate(func(rp *resolvers.StatusesRootResolver, rv resolvers.StatusesRootResolver) *resolvers.StatusesRootResolver {
		*rp = rv
		return &rv
	})
}
