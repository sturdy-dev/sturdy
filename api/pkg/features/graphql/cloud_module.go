//go:build cloud
// +build cloud

package graphql

import (
	"getsturdy.com/api/pkg/di"
	cloud "getsturdy.com/api/pkg/features/enterprise/cloud/graphql"
)

func Module(c *di.Container) {
	c.Import(cloud.Module)
}
