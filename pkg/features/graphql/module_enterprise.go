// +build enterprise

package graphql

import (
	"mash/pkg/di"
	"mash/pkg/features/graphql/enterprise"
)

func Module(c *di.Container) {
	c.Register(enterprise.NewFeaturesRootResolver)
}
