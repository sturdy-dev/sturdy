package graphql

import (
	"mash/pkg/di"
)

var Module = di.NewModule(
	di.ProvidesCycle(NewCodebaseRootResolver),
)
