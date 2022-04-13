package graphql

import (
	graphql_activity "getsturdy.com/api/pkg/activity/graphql"
	service_auth "getsturdy.com/api/pkg/auth/service"
	graphql_author "getsturdy.com/api/pkg/author/graphql"
	service_changes "getsturdy.com/api/pkg/changes/service"
	db_comments "getsturdy.com/api/pkg/comments/db"
	"getsturdy.com/api/pkg/di"
	graphql_downloads "getsturdy.com/api/pkg/downloads/graphql"
	graphql_file "getsturdy.com/api/pkg/file/graphql"
	"getsturdy.com/api/pkg/graphql/resolvers"
	"getsturdy.com/api/pkg/logger"
	"getsturdy.com/api/vcs/executor"
)

func Module(c *di.Container) {
	c.Import(service_changes.Module)
	c.Import(db_comments.Module)
	c.Import(service_auth.Module)
	c.Import(resolvers.Module)
	c.Import(graphql_author.Module)
	c.Import(graphql_downloads.Module)
	c.Import(graphql_activity.Module)
	c.Import(executor.Module)
	c.Import(logger.Module)
	c.Import(graphql_file.Module)
	c.Register(NewFileDiffRootResolver)
	c.Register(NewResolver)

	// populate cyclic resolver
	c.Decorate(func(rp *resolvers.ChangeRootResolver, rv resolvers.ChangeRootResolver) *resolvers.ChangeRootResolver {
		*rp = rv
		return &rv
	})
}
