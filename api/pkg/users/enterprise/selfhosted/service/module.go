package service

import (
	"getsturdy.com/api/pkg/di"
	selfhosted_service "getsturdy.com/api/pkg/users/oss/selfhosted/service"
	"getsturdy.com/api/pkg/users/service"
)

func Module(c *di.Container) {
	c.Import(selfhosted_service.Module)
	c.Register(New, new(service.Service))
}
