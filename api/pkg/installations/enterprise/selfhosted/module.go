package selfhosted

import (
	"getsturdy.com/api/pkg/di"
	"getsturdy.com/api/pkg/installations/enterprise/selfhosted/graphql"
	"getsturdy.com/api/pkg/installations/enterprise/selfhosted/worker"
)

func Module(c *di.Container) {
	c.Import(worker.Module)
	c.Import(graphql.Module)
}
