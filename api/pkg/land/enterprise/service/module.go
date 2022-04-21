package service

import (
	"getsturdy.com/api/pkg/di"
	service_github "getsturdy.com/api/pkg/github/enterprise/service"
	service_land "getsturdy.com/api/pkg/land/service"
	service_remote "getsturdy.com/api/pkg/remote/enterprise/service"
)

func Module(c *di.Container) {
	c.Import(service_github.Module)
	c.Import(service_remote.Module)
	c.Import(service_land.Module)
	c.Register(New)
}
