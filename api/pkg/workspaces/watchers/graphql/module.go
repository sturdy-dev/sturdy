package graphql

import (
	service_auth "getsturdy.com/api/pkg/auth/service"
	"getsturdy.com/api/pkg/di"
	"getsturdy.com/api/pkg/events/v2"
	"getsturdy.com/api/pkg/graphql/resolvers"
	"getsturdy.com/api/pkg/logger"
	service_workspace "getsturdy.com/api/pkg/workspaces/service"
	service_workspace_watchers "getsturdy.com/api/pkg/workspaces/watchers/service"
)

func Module(c *di.Container) {
	c.Import(logger.Module)
	c.Import(service_workspace.Module)
	c.Import(service_workspace_watchers.Module)
	c.Import(service_auth.Module)
	c.Import(events.Module)
	c.Import(resolvers.Module)
	c.Register(NewRootResolver)
}
