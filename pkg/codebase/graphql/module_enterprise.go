//go:build enterprise
// +build enterprise

package graphql

import (
	"mash/pkg/codebase/graphql/enterprise"
	"mash/pkg/codebase/graphql/oss"
	"mash/pkg/di"
	"mash/pkg/graphql/resolvers"
)

var Module = di.NewModule(
	di.Provides(oss.NewCodebaseRootResolver),
	di.Provides(enterprise.NewCodebaseRootResolver),
	di.ProvidesCycle(func(r *enterprise.CodebaseRootResolver) resolvers.CodebaseRootResolver {
		return r
	}),
)
