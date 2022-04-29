package graphql

import (
	service_analytics "getsturdy.com/api/pkg/analytics/service"
	service_auth "getsturdy.com/api/pkg/auth/service"
	grapqhl_author "getsturdy.com/api/pkg/author/graphql"
	"getsturdy.com/api/pkg/di"
	"getsturdy.com/api/pkg/events"
	eventsv2 "getsturdy.com/api/pkg/events/v2"
	"getsturdy.com/api/pkg/graphql/resolvers"
	"getsturdy.com/api/pkg/logger"
	"getsturdy.com/api/pkg/notification/sender"
	db_review "getsturdy.com/api/pkg/review/db"
	db_workspaces "getsturdy.com/api/pkg/workspaces/db"
	service_workspace_watchers "getsturdy.com/api/pkg/workspaces/watchers/service"
)

func Module(c *di.Container) {
	c.Import(logger.Module)
	c.Import(db_review.Module)
	c.Import(db_workspaces.Module)
	c.Import(service_auth.Module)
	c.Import(grapqhl_author.Module)
	c.Import(resolvers.Module)
	c.Import(events.Module)
	c.Import(eventsv2.Module)
	c.Import(sender.Module)
	c.Import(service_analytics.Module)
	c.Import(service_workspace_watchers.Module)
	c.Register(New)
}
