//go:build !enterprise
// +build !enterprise

package module

import (
	"mash/pkg/di"
	"mash/pkg/license/oss/graphql"
)

func Module(c *di.Container) {
	c.Register(graphql.New)
}
