//go:build !enterprise && !cloud
// +build !enterprise,!cloud

package module

import (
	"getsturdy.com/api/pkg/di"
	"getsturdy.com/api/pkg/workspace/service"
)

func Module(c *di.Container) {
	c.Register(service.New, new(service.Service))
}
