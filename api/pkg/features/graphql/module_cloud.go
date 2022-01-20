//go:build cloud
// +build cloud

package graphql

import (
	"getsturdy.com/api/pkg/di"
	"getsturdy.com/api/pkg/features/graphql/cloud"
)

func Module(c *di.Container) {
	c.Register(cloud.NewFeaturesRootResolver)
}
