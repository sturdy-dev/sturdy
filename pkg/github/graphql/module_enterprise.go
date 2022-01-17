//go:build enterprise || cloud
// +build enterprise cloud

package graphql

import (
	"mash/pkg/di"
	"mash/pkg/github/graphql/enterprise"
	"mash/pkg/graphql/resolvers"
)

func Module(c *di.Container) {
	c.Register(enterprise.NewGitHubAccountRootResolver, new(resolvers.GitHubAccountRootResolver))
	c.Register(enterprise.NewGitHubAppRootResolver)
	c.Register(enterprise.NewResolver)
}
