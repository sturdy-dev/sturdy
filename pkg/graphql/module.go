package graphql

import (
	"mash/pkg/di"
	graphql_user "mash/pkg/user/graphql"
	graphql_view "mash/pkg/view/graphql"
	graphql_workspace_activity "mash/pkg/workspace/activity/graphql"
	graphql_workspace "mash/pkg/workspace/graphql"
	graphql_workspace_watchers "mash/pkg/workspace/watchers/graphql"
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
