//go:build enterprise
// +build enterprise

package graphql

import (
	"mash/pkg/di"
	"mash/pkg/github/graphql/enterprise"
)

var Module = di.NewModule(
	di.Provides(enterprise.NewGitHubAppRootResolver),
	di.ProvidesCycle(enterprise.NewResolver),
)
