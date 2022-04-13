//go:build enterprise || cloud
// +build enterprise cloud

package graphql

import (
	"getsturdy.com/api/pkg/di"
	enterprise "getsturdy.com/api/pkg/github/enterprise/graphql"
)

func Module(c *di.Container) {
	c.Import(enterprise.Module)
}
