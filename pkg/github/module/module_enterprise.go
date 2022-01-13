//go:build enterprise
// +build enterprise

package module

import (
	"mash/pkg/di"
	"mash/pkg/github/enterprise"
	"mash/pkg/github/graphql"
	graphql_pr "mash/pkg/github/graphql/pr"
)

func Module(c *di.Container) {
	c.Import(enterprise.Module)
	c.Import(graphql.Module)
	c.Import(graphql_pr.Module)
}
