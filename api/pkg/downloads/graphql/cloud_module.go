//go:build cloud
// +build cloud

package graphql

import (
	"getsturdy.com/api/pkg/di"
	"getsturdy.com/api/pkg/downloads/enterprise/cloud/graphql"
)

func Module(c *di.Container) {
	c.Import(graphql.Module)
}
