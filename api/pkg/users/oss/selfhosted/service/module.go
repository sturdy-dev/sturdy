package service

import (
	"getsturdy.com/api/pkg/di"
	service_organization "getsturdy.com/api/pkg/organization/service"
	service_users "getsturdy.com/api/pkg/users/service"
)

func Module(c *di.Container) {
	c.Import(service_users.Module)
	c.Import(service_organization.Module)
	c.Register(New)
}
