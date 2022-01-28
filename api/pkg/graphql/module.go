package graphql

import (
	"getsturdy.com/api/pkg/di"
	graphql_workspace_activity "getsturdy.com/api/pkg/workspace/activity/graphql"
	graphql_workspace_watchers "getsturdy.com/api/pkg/workspace/watchers/graphql"
)

func Module(c *di.Container) {
	c.Register(NewRootResolver)

	c.Register(graphql_workspace_watchers.NewRootResolver)
	c.Register(graphql_workspace_activity.New)
}
