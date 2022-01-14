//go:build !enterprise
// +build !enterprise

package graphql

import (
	"mash/pkg/di"
	"mash/pkg/features/graphql/oss"
)

func Module(c *di.Container) {
	c.Register(oss.NewFeaturesRootResolver)
}
