package service

import (
	service_codebases "getsturdy.com/api/pkg/codebases/service"
	"getsturdy.com/api/pkg/di"
	service_installations "getsturdy.com/api/pkg/installations/service"
	"getsturdy.com/api/pkg/installations/statistics/enterprise/selfhosted/publisher"
	service_users "getsturdy.com/api/pkg/users/service/module"
)

func Module(c *di.Container) {
	c.Import(service_installations.Module)
	c.Import(service_codebases.Module)
	c.Import(service_users.Module)
	c.Import(publisher.Module)
	c.Register(New)
}
