package service

import (
	"getsturdy.com/api/pkg/di"
	"getsturdy.com/api/pkg/installations/db"
	"getsturdy.com/api/pkg/installations/service"
	service_statistics "getsturdy.com/api/pkg/installations/statistics/enterprise/selfhosted/service"
	"getsturdy.com/api/pkg/licenses/enterprise/selfhosted/validator"
)

func Module(c *di.Container) {
	c.Import(service.Module)
	c.Import(validator.Module)
	c.Import(db.Module)
	c.Import(service_statistics.Module)
	c.Register(New)
}
