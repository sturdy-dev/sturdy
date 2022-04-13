package graphql

import (
	db_activity "getsturdy.com/api/pkg/activity/db"
	service_activity "getsturdy.com/api/pkg/activity/service"
	service_auth "getsturdy.com/api/pkg/auth/service"
	graphql_authors "getsturdy.com/api/pkg/author/graphql"
	"getsturdy.com/api/pkg/di"
	"getsturdy.com/api/pkg/graphql/resolvers"
	graphql_review "getsturdy.com/api/pkg/review/graphql"
)

func Module(c *di.Container) {
	c.Import(graphql_authors.Module)
	c.Import(db_activity.Module)
	c.Import(graphql_review.Module)
	c.Import(resolvers.Module)
	c.Import(service_activity.Module)
	c.Import(service_auth.Module)
	c.Register(New)
}
