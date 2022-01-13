//go:build enterprise
// +build enterprise

package graphql

import (
	"mash/pkg/di"
	"mash/pkg/graphql/resolvers"
	"mash/pkg/statuses/graphql"
	"mash/pkg/statuses/graphql/enterprise"
)

func Module(c *di.Container) {
	c.Register(graphql.New)
	c.Register(enterprise.New, new(resolvers.StatusesRootResolver))
}
