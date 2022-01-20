//go:build !enterprise && !cloud
// +build !enterprise,!cloud

package graphql

import (
	"getsturdy.com/api/pkg/di"
	"getsturdy.com/api/pkg/graphql/resolvers"
	"getsturdy.com/api/pkg/statuses/graphql"
)

func Module(c *di.Container) {
	c.Register(graphql.New, new(resolvers.StatusesRootResolver))
}
