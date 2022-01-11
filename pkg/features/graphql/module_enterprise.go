// +build enterprise

package grapqhl

import (
	"mash/pkg/di"
	"mash/pkg/features/graphql/enterprise"
)

var Module = di.NewModule(
	di.Provides(enterprise.NewFeaturesRootResolver),
)
