package graphql

import (
	"getsturdy.com/api/pkg/di"
	"getsturdy.com/api/pkg/graphql/resolvers"
)

func Module(c *di.Container) {
	c.Register(NewGitHubAccountRootResolver, new(resolvers.GitHubAccountRootResolver))
	c.Register(NewGitHubAppRootResolver)
	c.Register(NewCodebaseGitHubIntegrationRootResolver)
	c.Register(NewGitHubRootResolver)
}
