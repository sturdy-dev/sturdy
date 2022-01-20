//go:build !enterprise && !cloud
// +build !enterprise,!cloud

package module

import (
	"getsturdy.com/api/pkg/di"
	"getsturdy.com/api/pkg/github/graphql"
	graphql_pr "getsturdy.com/api/pkg/github/graphql/pr"
)

func Module(c *di.Container) {
	c.Import(graphql.Module)
	c.Import(graphql_pr.Module)
}
