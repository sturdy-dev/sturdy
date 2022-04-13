//go:build !enterprise && !cloud
// +build !enterprise,!cloud

package module

import (
	"getsturdy.com/api/pkg/di"
	oss "getsturdy.com/api/pkg/users/oss/selfhosted/service"
	"getsturdy.com/api/pkg/users/service"
)

func Module(c *di.Container) {
	c.Import(oss.Module)
	c.Register(func(s *oss.Service) service.Service {
		return s
	})
}
