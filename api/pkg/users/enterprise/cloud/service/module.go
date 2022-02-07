package service

import (
	"getsturdy.com/api/pkg/di"
	service_users "getsturdy.com/api/pkg/users/service"
)

func Module(c *di.Container) {
	c.Register(New)
	c.Register(func(service *Service) service_users.Service {
		return service
	})
}
