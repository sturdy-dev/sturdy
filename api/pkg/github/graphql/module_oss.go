//go:build !enterprise && !cloud
// +build !enterprise,!cloud

package graphql

import (
	"getsturdy.com/api/pkg/di"
	"getsturdy.com/api/pkg/github/graphql/oss"
	"getsturdy.com/api/pkg/graphql/resolvers"
)

func Module(c *di.Container) {
	c.Register(oss.NewGitHubAccountRootResolver, new(resolvers.GitHubAccountRootResolver))
	c.Register(oss.NewGitHubAppRootResolver)
	c.Register(oss.NewCodebaseGitHubIntegrationRootResolver)
	c.Register(oss.NewGitHubRootResolver)
}
