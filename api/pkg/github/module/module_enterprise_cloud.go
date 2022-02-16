//go:build enterprise || cloud
// +build enterprise cloud

package module

import (
	"getsturdy.com/api/pkg/di"
	gh_enterprise "getsturdy.com/api/pkg/github/enterprise"
	graphql_github "getsturdy.com/api/pkg/github/graphql"
	graphql_pr "getsturdy.com/api/pkg/github/graphql/pr"
)

func Module(c *di.Container) {
	c.Import(gh_enterprise.Module)
	c.Import(graphql_github.Module)
	c.Import(graphql_pr.Module)
}
