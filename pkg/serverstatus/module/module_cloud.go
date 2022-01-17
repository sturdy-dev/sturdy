//go:build cloud
// +build cloud

package module

import (
	"mash/pkg/di"
	"mash/pkg/serverstatus/cloud/graphql"
)

func Module(c *di.Container) {
	c.Register(graphql.New)
}
