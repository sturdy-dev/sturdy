//go:build !enterprise
// +build !enterprise

package graphql

import (
	"mash/pkg/di"
	"mash/pkg/github/graphql/oss"
)

var Module = di.NewModule(
	di.Provides(oss.NewGitHubAppRootResolver),
	di.ProvidesCycle(oss.NewResolver),
)
