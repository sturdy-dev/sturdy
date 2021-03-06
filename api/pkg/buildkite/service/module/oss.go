//go:build !enterprise && !cloud
// +build !enterprise,!cloud

package module

import (
	"getsturdy.com/api/pkg/buildkite/service"
	"getsturdy.com/api/pkg/di"
)

func Module(c *di.Container) {
	c.Register(service.New, new(service.Service))
}
