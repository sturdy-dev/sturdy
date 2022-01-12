//go:build !enterprise
// +build !enterprise

package module

import (
	"mash/pkg/di"
	"mash/pkg/license/oss/graphql"
)

var Module = di.NewModule(
	di.ProvidesCycle(graphql.New),
)
