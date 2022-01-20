//go:build cloud
// +build cloud

package module

import (
	"getsturdy.com/api/pkg/di"
	"getsturdy.com/api/pkg/serverstatus/cloud/graphql"
)

func Module(c *di.Container) {
	c.Register(graphql.New)
}
