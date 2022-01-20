//go:build enterprise || cloud
// +build enterprise cloud

package graphql

import (
	"getsturdy.com/api/pkg/di"
	"getsturdy.com/api/pkg/github/graphql/enterprise"
	"getsturdy.com/api/pkg/graphql/resolvers"
)

func Module(c *di.Container) {
	c.Register(enterprise.NewGitHubAccountRootResolver, new(resolvers.GitHubAccountRootResolver))
	c.Register(enterprise.NewGitHubAppRootResolver)
	c.Register(enterprise.NewResolver)
}
