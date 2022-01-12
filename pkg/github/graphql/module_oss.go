//go:build !enterprise
// +build !enterprise

package graphql

import (
	"mash/pkg/di"
	"mash/pkg/github/graphql/oss"
	"mash/pkg/graphql/resolvers"
)

var Module = di.NewModule(
	di.Provides(oss.NewGitHubAccountRootResolver, new(resolvers.GitHubAccountRootResolver)),
	di.Provides(oss.NewGitHubAppRootResolver),
	di.ProvidesCycle(oss.NewResolver),
)
