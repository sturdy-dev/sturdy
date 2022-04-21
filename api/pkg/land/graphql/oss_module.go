//go:build !cloud && !enterprise
// +build !cloud,!enterprise

package grapqhl

import (
	services_auth "getsturdy.com/api/pkg/auth/service"
	"getsturdy.com/api/pkg/di"
	service_land "getsturdy.com/api/pkg/land/service"
	graphql_workspaces "getsturdy.com/api/pkg/workspaces/graphql"
	service_workspaces "getsturdy.com/api/pkg/workspaces/service"
)

func Module(c *di.Container) {
	c.Import(service_land.Module)
	c.Import(services_auth.Module)
	c.Import(service_workspaces.Module)
	c.Import(graphql_workspaces.Module)
	c.Register(NewResolver)
}
