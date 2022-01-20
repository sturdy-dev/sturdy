//go:build !cloud
// +build !cloud

package module

import (
	"getsturdy.com/api/pkg/serverstatus/selfhosted/graphql"
	"getsturdy.com/api/pkg/serverstatus/selfhosted/service"

	"getsturdy.com/api/pkg/di"
)

func Module(c *di.Container) {
	c.Register(graphql.New)
	c.Register(service.New)
}
