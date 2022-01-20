//go:build enterprise || cloud
// +build enterprise cloud

package graphql

import (
	"getsturdy.com/api/pkg/di"
	"getsturdy.com/api/pkg/graphql/resolvers"
	"getsturdy.com/api/pkg/statuses/graphql"
	"getsturdy.com/api/pkg/statuses/graphql/enterprise"
)

func Module(c *di.Container) {
	c.Register(graphql.New)
	c.Register(enterprise.New, new(resolvers.StatusesRootResolver))
}
