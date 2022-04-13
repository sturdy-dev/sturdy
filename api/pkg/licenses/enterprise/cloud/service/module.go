package service

import (
	"getsturdy.com/api/pkg/di"
	service_installation_statistics "getsturdy.com/api/pkg/installations/statistics/enterprise/cloud/service"
	"getsturdy.com/api/pkg/licenses/enterprise/cloud/db"
	service_license_validations "getsturdy.com/api/pkg/licenses/enterprise/cloud/validations/service"
)

func Module(c *di.Container) {
	c.Import(db.Module)
	c.Import(service_installation_statistics.Module)
	c.Import(service_license_validations.Module)
	c.Register(NewService)
}
