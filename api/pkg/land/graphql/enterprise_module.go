//go:build cloud || enterprise
// +build cloud enterprise

package grapqhl

import (
	"getsturdy.com/api/pkg/di"
	graphql_land "getsturdy.com/api/pkg/land/enterprise/graphql"
)

func Module(c *di.Container) {
	c.Import(graphql_land.Module)
}
