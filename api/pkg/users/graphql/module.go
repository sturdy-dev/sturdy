package graphql

import (
	service_analytics "getsturdy.com/api/pkg/analytics/service"
	"getsturdy.com/api/pkg/di"
	graphql_github "getsturdy.com/api/pkg/github/graphql"
	"getsturdy.com/api/pkg/graphql/resolvers"
	"getsturdy.com/api/pkg/logger"
	db_newsletter "getsturdy.com/api/pkg/newsletter/db"
	graphql_notification "getsturdy.com/api/pkg/notification/graphql"
	db_user "getsturdy.com/api/pkg/users/db"
	service_user "getsturdy.com/api/pkg/users/service"
	graphql_view "getsturdy.com/api/pkg/views/graphql"
)

func Module(c *di.Container) {
	c.Import(db_user.Module)
	c.Import(db_newsletter.Module)
	c.Import(service_user.Module)
	c.Import(graphql_view.Module)
	c.Import(graphql_notification.Module)
	c.Import(graphql_github.Module)
	c.Import(logger.Module)
	c.Import(service_analytics.Module)
	c.Register(NewResolver, new(resolvers.UserRootResolver))

	// populate cyclic resolver
	c.Decorate(func(rp *resolvers.UserRootResolver, rv resolvers.UserRootResolver) *resolvers.UserRootResolver {
		*rp = rv
		return &rv
	})
}
