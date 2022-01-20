package module

import (
	"getsturdy.com/api/pkg/di"
	"getsturdy.com/api/pkg/gc/db"
	"getsturdy.com/api/pkg/gc/worker"
)

func Module(c *di.Container) {
	c.Import(db.Module)
	c.Import(worker.Module)
}
