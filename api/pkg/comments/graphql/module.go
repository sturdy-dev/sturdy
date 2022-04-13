package graphql

import (
	sender_workspace_activity "getsturdy.com/api/pkg/activity/sender"
	service_analytics "getsturdy.com/api/pkg/analytics/service"
	service_auth "getsturdy.com/api/pkg/auth/service"
	graphql_author "getsturdy.com/api/pkg/author/graphql"
	graphql_changes "getsturdy.com/api/pkg/changes/graphql"
	service_change "getsturdy.com/api/pkg/changes/service"
	db_codebases "getsturdy.com/api/pkg/codebases/db"
	db_comments "getsturdy.com/api/pkg/comments/db"
	"getsturdy.com/api/pkg/di"
	"getsturdy.com/api/pkg/events"
	eventsv2 "getsturdy.com/api/pkg/events/v2"
	"getsturdy.com/api/pkg/graphql/resolvers"
	"getsturdy.com/api/pkg/logger"
	notification_sender "getsturdy.com/api/pkg/notification/sender"
	db_snapshots "getsturdy.com/api/pkg/snapshots/db"
	service_users "getsturdy.com/api/pkg/users/service/module"
	db_view "getsturdy.com/api/pkg/view/db"
	db_workspaces "getsturdy.com/api/pkg/workspaces/db"
	service_workspace_watchers "getsturdy.com/api/pkg/workspaces/watchers/service"
	"getsturdy.com/api/vcs/executor"
)

func Module(c *di.Container) {
	c.Import(db_comments.Module)
	c.Import(db_snapshots.Module)
	c.Import(db_workspaces.Module)
	c.Import(db_view.Module)
	c.Import(db_codebases.Module)
	c.Import(service_workspace_watchers.Module)
	c.Import(service_auth.Module)
	c.Import(service_change.Module)
	c.Import(events.Module)
	c.Import(eventsv2.Module)
	c.Import(notification_sender.Module)
	c.Import(sender_workspace_activity.Module)
	c.Import(service_users.Module)
	c.Import(graphql_author.Module)
	c.Import(graphql_changes.Module)
	c.Import(resolvers.Module)
	c.Import(logger.Module)
	c.Import(service_analytics.Module)
	c.Import(executor.Module)
	c.Register(NewResolver)

	// populate cyclic resolver
	c.Decorate(func(rp *resolvers.CommentRootResolver, rv resolvers.CommentRootResolver) *resolvers.CommentRootResolver {
		*rp = rv
		return &rv
	})
}
