//go:build !enterprise
// +build !enterprise

package graphql

import (
	"mash/pkg/codebase/graphql/oss"
	"mash/pkg/di"
	"mash/pkg/graphql/resolvers"
)

var Module = di.NewModule(
	di.Provides(oss.NewCodebaseRootResolver),
	di.ProvidesCycle(func(r *oss.CodebaseRootResolver) resolvers.CodebaseRootResolver {
		return r
	}),
)
