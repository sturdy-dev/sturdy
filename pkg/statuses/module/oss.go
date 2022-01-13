//go:build !enterprise
// +build !enterprise

package graphql

import (
	"mash/pkg/di"
	"mash/pkg/graphql/resolvers"
	"mash/pkg/statuses/graphql"
)

func Module(c *di.Container) {
	c.Register(graphql.New, new(resolvers.StatusesRootResolver))
}
