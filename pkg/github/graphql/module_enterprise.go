//go:build enterprise
// +build enterprise

package graphql

import (
	"mash/pkg/di"
	"mash/pkg/github/graphql/enterprise"
	"mash/pkg/graphql/resolvers"
)

var Module = di.NewModule(
	di.Provides(enterprise.NewGitHubAccountRootResolver, new(resolvers.GitHubAccountRootResolver)),
	di.Provides(enterprise.NewGitHubAppRootResolver),
	di.ProvidesCycle(enterprise.NewResolver),
)
