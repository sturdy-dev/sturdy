package cloud

import (
	"getsturdy.com/api/pkg/di"
	"getsturdy.com/api/pkg/installations/statistics/enterprise/cloud/db"
	"getsturdy.com/api/pkg/installations/statistics/enterprise/cloud/service"
)

func Module(c *di.Container) {
	c.Import(service.Module)
	c.Import(db.Module)
}
