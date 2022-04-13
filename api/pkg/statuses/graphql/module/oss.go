//go:build !enterprise && !cloud
// +build !enterprise,!cloud

package graphql

import (
	"getsturdy.com/api/pkg/di"
	"getsturdy.com/api/pkg/statuses/graphql"
)

func Module(c *di.Container) {
	c.Import(graphql.Module)
}
