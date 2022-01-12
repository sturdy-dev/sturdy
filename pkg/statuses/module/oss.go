//go:build !enterprise
// +build !enterprise

package graphql

import (
	"mash/pkg/di"
	"mash/pkg/graphql/resolvers"
	"mash/pkg/statuses/graphql"
)

var Module = di.NewModule(
	di.Provides(graphql.New),
	di.ProvidesCycle(func(r *graphql.RootResolver) resolvers.StatusesRootResolver {
		return r
	}),
)
