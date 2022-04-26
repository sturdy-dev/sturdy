//go:build enterprise || cloud
// +build enterprise cloud

package module

import (
	"getsturdy.com/api/pkg/di"
	service_github "getsturdy.com/api/pkg/github/enterprise/service"
	"getsturdy.com/api/pkg/github/service"
)

func Module(c *di.Container) {
	c.Import(service_github.Module)
	c.Register(func(svc *service_github.Service) *service_github.Service { return svc }, new(service.Service))
}
