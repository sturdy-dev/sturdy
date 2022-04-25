//go:build !cloud && !enterprise

package graphql

import (
	"getsturdy.com/api/pkg/buildkite/graphql"
	"getsturdy.com/api/pkg/di"
)

func Module(c *di.Container) {
	c.Register(graphql.New)
}
