//go:build enterprise
// +build enterprise

package graphql

import (
	"mash/pkg/di"
	"mash/pkg/graphql/resolvers"
	"mash/pkg/statuses/graphql"
	"mash/pkg/statuses/graphql/enterprise"
)

var Module = di.NewModule(
	di.Provides(graphql.New),
	di.Provides(enterprise.New),
	di.ProvidesCycle(func(r *enterprise.RootResolver) resolvers.StatusesRootResolver {
		return r
	}),
)
