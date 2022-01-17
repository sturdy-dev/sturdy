//go:build enterprise || cloud
// +build enterprise cloud

package module

import (
	"mash/pkg/di"
	"mash/pkg/workspace/service"
	"mash/pkg/workspace/service/enterprise"
)

func Module(c *di.Container) {
	c.Register(service.New)
	c.Register(enterprise.New, new(service.Service))
}
