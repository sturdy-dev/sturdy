package service

import (
	"getsturdy.com/api/pkg/di"
	service_github "getsturdy.com/api/pkg/github/enterprise/service"
	service_remote "getsturdy.com/api/pkg/remote/enterprise/service"
	"getsturdy.com/api/pkg/workspaces/service"
)

func Module(c *di.Container) {
	c.Import(service.Module)
	c.Import(service_github.Module)
	c.Import(service_remote.Module)
	c.Register(New)
}
