package selfhosted

import (
	"getsturdy.com/api/pkg/api"
	"getsturdy.com/api/pkg/di"
	workers_license "getsturdy.com/api/pkg/installations/enterprise/selfhosted/worker"
	worker_installation_statistics "getsturdy.com/api/pkg/installations/statistics/enterprise/selfhosted/worker"
)

func Module(c *di.Container) {
	c.Import(api.Module)
	c.Import(workers_license.Module)
	c.Import(worker_installation_statistics.Module)
	c.Register(ProvideAPI, new(api.Starter))
}
