//go:build enterprise
// +build enterprise

package graphql

import (
	"getsturdy.com/api/pkg/di"
	"getsturdy.com/api/pkg/features/graphql/enterprise"
)

func Module(c *di.Container) {
	c.Register(enterprise.NewFeaturesRootResolver)
}
