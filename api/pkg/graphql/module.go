package graphql

import (
	"getsturdy.com/api/pkg/di"
	graphql_user "getsturdy.com/api/pkg/user/graphql"
	graphql_view "getsturdy.com/api/pkg/view/graphql"
	graphql_workspace_activity "getsturdy.com/api/pkg/workspace/activity/graphql"
	graphql_workspace "getsturdy.com/api/pkg/workspace/graphql"
	graphql_workspace_watchers "getsturdy.com/api/pkg/workspace/watchers/graphql"
)

func Module(c *di.Container) {
	c.Register(NewRootResolver)

	c.Register(graphql_user.NewResolver)
	c.Register(graphql_view.NewViewStatusRootResolver)
	c.Register(graphql_view.NewResolver)
	c.Register(graphql_workspace_watchers.NewRootResolver)
	c.Register(graphql_workspace.NewResolver)
	c.Register(graphql_workspace_activity.New)
}
