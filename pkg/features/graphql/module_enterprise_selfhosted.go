//go:build enterprise && !cloud
// +build enterprise,!cloud

package grapqhl

import (
	"mash/pkg/di"
	"mash/pkg/features/graphql/enterprise/selfhosted"
)

var Module = di.NewModule(
	di.Provides(selfhosted.NewFeaturesRootResolver),
)
