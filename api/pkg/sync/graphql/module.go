package graphql

import (
	graphql_changes "getsturdy.com/api/pkg/changes/graphql"
	"getsturdy.com/api/pkg/di"
	"getsturdy.com/api/pkg/logger"
	service_workspace "getsturdy.com/api/pkg/workspaces/service"
	"getsturdy.com/api/vcs/executor"
)

func Module(c *di.Container) {
	c.Import(logger.Module)
	c.Import(executor.Module)
	c.Import(graphql_changes.Module)
	c.Import(service_workspace.Module)
	c.Register(NewRootResolver)
}
