package graphql

import (
	"getsturdy.com/api/pkg/di"
	"getsturdy.com/api/pkg/graphql/resolvers"
	"getsturdy.com/api/pkg/logger"
	service_users "getsturdy.com/api/pkg/users/service/module"
)

func Module(c *di.Container) {
	c.Import(logger.Module)
	c.Import(service_users.Module)
	c.Register(NewResolver)

	// populate cyclic resolver
	c.Import(resolvers.Module)
	c.Decorate(func(rp *resolvers.AuthorRootResolver, rv resolvers.AuthorRootResolver) *resolvers.AuthorRootResolver {
		*rp = rv
		return &rv
	})
}
