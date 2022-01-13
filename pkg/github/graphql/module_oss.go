//go:build !enterprise
// +build !enterprise

package graphql

import (
	"mash/pkg/di"
	"mash/pkg/github/graphql/oss"
	"mash/pkg/graphql/resolvers"
)

func Module(c *di.Container) {
	c.Register(oss.NewGitHubAccountRootResolver, new(resolvers.GitHubAccountRootResolver))
	c.Register(oss.NewGitHubAppRootResolver)
	c.Register(oss.NewResolver)
}
