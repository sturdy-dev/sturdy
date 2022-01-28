package selfhosted

import (
	"getsturdy.com/api/pkg/di"
	"getsturdy.com/api/pkg/installations/statistics/enterprise/selfhosted/publisher"
	"getsturdy.com/api/pkg/installations/statistics/enterprise/selfhosted/service"
	"getsturdy.com/api/pkg/installations/statistics/enterprise/selfhosted/worker"
)

func Module(c *di.Container) {
	c.Import(worker.Module)
	c.Import(publisher.Module)
	c.Import(service.Module)
}
