// +build enterprise

package grapqhl

import (
	"mash/pkg/di"
	"mash/pkg/features/graphql/enterprise"
)

func Module(c *di.Container) {
	c.Register(enterprise.NewFeaturesRootResolver)
}
