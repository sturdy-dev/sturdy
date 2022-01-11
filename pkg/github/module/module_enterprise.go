//go:build enterprise
// +build enterprise

package module

import (
	"mash/pkg/di"
	"mash/pkg/github/enterprise"
	"mash/pkg/github/graphql"
	graphql_pr "mash/pkg/github/graphql/pr"
)

var Module = di.NewModule(
	di.Needs(enterprise.Module),
	di.Needs(graphql.Module),
	di.Needs(graphql_pr.Module),
)
