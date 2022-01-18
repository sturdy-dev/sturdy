//go:build !enterprise && !cloud
// +build !enterprise,!cloud

package module

import (
	"mash/pkg/di"
	"mash/pkg/github/graphql"
	graphql_pr "mash/pkg/github/graphql/pr"
)

func Module(c *di.Container) {
	c.Import(graphql.Module)
	c.Import(graphql_pr.Module)
}
