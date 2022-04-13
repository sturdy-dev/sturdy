//go:build !cloud && !enterprise
// +build !cloud,!enterprise

package module

import (
	"getsturdy.com/api/pkg/di"
	"getsturdy.com/api/pkg/workspaces/service"
)

func Module(c *di.Container) {
	c.Import(service.Module)
	c.Register(func(s *service.WorkspaceService) service.Service {
		return s
	})
}
