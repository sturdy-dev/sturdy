//go:build !enterprise
// +build !enterprise

package pr

import (
	"mash/pkg/di"
	"mash/pkg/github/graphql/pr/oss"
)

var Module = di.NewModule(
	di.ProvidesCycle(oss.NewResolver),
)
