//go:build !enterprise && !cloud
// +build !enterprise,!cloud

package graphql

import (
	"getsturdy.com/api/pkg/di"
	"getsturdy.com/api/pkg/graphql/resolvers"
	"getsturdy.com/api/pkg/statuses/db"
	"getsturdy.com/api/pkg/statuses/graphql"
	"getsturdy.com/api/pkg/statuses/service"
)

func Module(c *di.Container) {
	c.Import(db.Module)
	c.Import(service.Module)
	c.Register(graphql.New, new(resolvers.StatusesRootResolver))
}
