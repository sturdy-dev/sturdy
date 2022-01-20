//go:build !enterprise && !cloud
// +build !enterprise,!cloud

package graphql

import (
	"getsturdy.com/api/pkg/di"
	"getsturdy.com/api/pkg/features/graphql/oss"
)

func Module(c *di.Container) {
	c.Register(oss.NewFeaturesRootResolver)
}
