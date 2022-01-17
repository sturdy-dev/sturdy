//go:build !cloud
// +build !cloud

package module

import (
	"mash/pkg/serverstatus/selfhosted/graphql"
	"mash/pkg/serverstatus/selfhosted/service"

	"mash/pkg/di"
)

func Module(c *di.Container) {
	c.Register(graphql.New)
	c.Register(service.New)
}
