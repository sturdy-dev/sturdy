package module

import (
	"getsturdy.com/api/pkg/di"
	"getsturdy.com/api/pkg/onetime/db"
	"getsturdy.com/api/pkg/onetime/service"
)

func Module(c *di.Container) {
	c.Import(db.Module)
	c.Import(service.Module)
}
