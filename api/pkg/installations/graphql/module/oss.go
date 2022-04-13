//go:build !enterprise && !cloud
// +build !enterprise,!cloud

package module

import (
	"getsturdy.com/api/pkg/di"
	"getsturdy.com/api/pkg/graphql/resolvers"
	"getsturdy.com/api/pkg/installations/graphql"
)

func Module(c *di.Container) {
	c.Import(graphql.Module)
	c.Register(func(r *graphql.RootResolver) resolvers.InstallationsRootResolver {
		return r
	})
}
