//go:build enterprise
// +build enterprise

package pr

import (
	"mash/pkg/di"
	"mash/pkg/github/graphql/pr/enterprise"
)

var Module = di.NewModule(
	di.ProvidesCycle(enterprise.NewResolver),
)
