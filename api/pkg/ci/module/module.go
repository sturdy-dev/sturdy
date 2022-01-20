package module

import (
	"getsturdy.com/api/pkg/ci/db"
	"getsturdy.com/api/pkg/ci/service"
	"getsturdy.com/api/pkg/ci/workers"
	"getsturdy.com/api/pkg/di"
)

func Module(c *di.Container) {
	c.Import(db.Module)
	c.Import(service.Module)
	c.Import(workers.Module)
}
