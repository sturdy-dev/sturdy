package graphql

import (
	service_analytics "getsturdy.com/api/pkg/analytics/service"
	service_auth "getsturdy.com/api/pkg/auth/service"
	graphql_author "getsturdy.com/api/pkg/author/graphql"
	graphql_codebases "getsturdy.com/api/pkg/codebases/graphql"
	"getsturdy.com/api/pkg/di"
	"getsturdy.com/api/pkg/events"
	eventsv2 "getsturdy.com/api/pkg/events/v2"
	"getsturdy.com/api/pkg/graphql/resolvers"
	"getsturdy.com/api/pkg/logger"
	db_mutagen "getsturdy.com/api/pkg/mutagen/db"
	service_snapshots "getsturdy.com/api/pkg/snapshots/service"
	db_view "getsturdy.com/api/pkg/views/db"
	service_view "getsturdy.com/api/pkg/views/service"
	db_workspaces "getsturdy.com/api/pkg/workspaces/db"
	meta_workspaces "getsturdy.com/api/pkg/workspaces/meta"
	service_workspace_watchers "getsturdy.com/api/pkg/workspaces/watchers/service"
	"getsturdy.com/api/vcs/executor"
)

func Module(c *di.Container) {
	c.Import(db_view.Module)
	c.Import(db_workspaces.Module)
	c.Import(service_snapshots.Module)
	c.Import(graphql_author.Module)
	c.Import(resolvers.Module)
	c.Import(meta_workspaces.Module)
	c.Import(events.Module)
	c.Import(eventsv2.Module)
	c.Import(executor.Module)
	c.Import(logger.Module)
	c.Import(service_workspace_watchers.Module)
	c.Import(service_analytics.Module)
	c.Import(graphql_codebases.Module)
	c.Import(service_auth.Module)
	c.Import(service_view.Module)
	c.Import(db_mutagen.Module)
	c.Register(NewResolver)
	c.Register(NewViewStatusRootResolver)

	// populate cyclic resolver
	c.Decorate(func(rp *resolvers.ViewRootResolver, rv resolvers.ViewRootResolver) *resolvers.ViewRootResolver {
		*rp = rv
		return &rv
	})
}
