//go:build !enterprise && !cloud
// +build !enterprise,!cloud

package graphql

import (
	"mash/pkg/di"
	"mash/pkg/features/graphql/oss"
)

func Module(c *di.Container) {
	c.Register(oss.NewFeaturesRootResolver)
}
