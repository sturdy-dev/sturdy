package worker

import (
	"getsturdy.com/api/pkg/di"
	service_installations "getsturdy.com/api/pkg/installations/service"
	validator_license "getsturdy.com/api/pkg/licenses/enterprise/selfhosted/validator"
	"getsturdy.com/api/pkg/logger"
)

func Module(c *di.Container) {
	c.Import(logger.Module)
	c.Import(validator_license.Module)
	c.Import(service_installations.Module)
	c.Register(New)
}
