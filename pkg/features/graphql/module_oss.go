// +build !enterprise

package graphql

import (
	"mash/pkg/di"
	"mash/pkg/features/graphql/oss"
)

var Module = di.NewModule(
	di.Provides(oss.NewFeaturesRootResolver),
)
