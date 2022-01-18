//go:build cloud
// +build cloud

package graphql

import (
	"mash/pkg/di"
	"mash/pkg/features/graphql/cloud"
)

func Module(c *di.Container) {
	c.Register(cloud.NewFeaturesRootResolver)
}
