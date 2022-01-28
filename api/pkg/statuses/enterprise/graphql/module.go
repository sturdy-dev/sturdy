package graphql

import (
	"getsturdy.com/api/pkg/di"
	"getsturdy.com/api/pkg/graphql/resolvers"
)

func Module(c *di.Container) {
	c.Register(New, new(resolvers.StatusesRootResolver))
}
