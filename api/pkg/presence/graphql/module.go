package graphql

import (
	graphql_author "getsturdy.com/api/pkg/author/graphql"
	"getsturdy.com/api/pkg/di"
	"getsturdy.com/api/pkg/events/v2"
	"getsturdy.com/api/pkg/graphql/resolvers"
	"getsturdy.com/api/pkg/logger"
	service_presence "getsturdy.com/api/pkg/presence/service"
)

func Module(c *di.Container) {
	c.Import(logger.Module)
	c.Import(service_presence.Module)
	c.Import(events.Module)
	c.Import(resolvers.Module)
	c.Import(graphql_author.Module)
	c.Register(NewRootResolver)
}
