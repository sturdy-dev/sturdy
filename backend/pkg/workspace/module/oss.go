//go:build !enterprise && !cloud
// +build !enterprise,!cloud

package module

import (
	"mash/pkg/di"
	"mash/pkg/workspace/service"
)

func Module(c *di.Container) {
	c.Register(service.New, new(service.Service))
}
