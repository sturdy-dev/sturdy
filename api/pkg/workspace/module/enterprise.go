//go:build enterprise || cloud
// +build enterprise cloud

package module

import (
	"getsturdy.com/api/pkg/di"
	"getsturdy.com/api/pkg/workspace/service"
	"getsturdy.com/api/pkg/workspace/service/enterprise"
)

func Module(c *di.Container) {
	c.Register(service.New)
	c.Register(enterprise.New, new(service.Service))
}
