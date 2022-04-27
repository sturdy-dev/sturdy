package schema

import (
	"getsturdy.com/api/pkg/db"
	"getsturdy.com/api/pkg/di"
	"getsturdy.com/api/pkg/logger"
)

func Module(c *di.Container) {
	c.Import(db.Module)
	c.Import(logger.Module)
	c.Register(New)
}
