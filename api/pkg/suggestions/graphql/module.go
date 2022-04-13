package graphql

import (
	service_auth "getsturdy.com/api/pkg/auth/service"
	graphql_author "getsturdy.com/api/pkg/author/graphql"
	graphql_changes "getsturdy.com/api/pkg/changes/graphql"
	"getsturdy.com/api/pkg/di"
	"getsturdy.com/api/pkg/events"
	"getsturdy.com/api/pkg/graphql/resolvers"
	"getsturdy.com/api/pkg/logger"
	service_suggestions "getsturdy.com/api/pkg/suggestions/service"
	service_workspace "getsturdy.com/api/pkg/workspaces/service"
)

func Module(c *di.Container) {
	c.Import(logger.Module)
	c.Import(service_auth.Module)
	c.Import(service_suggestions.Module)
	c.Import(service_workspace.Module)
	c.Import(graphql_author.Module)
	c.Import(graphql_changes.Module)
	c.Import(resolvers.Module)
	c.Import(events.Module)
	c.Register(New)
}
