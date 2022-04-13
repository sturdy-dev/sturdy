package service

import (
	"getsturdy.com/api/pkg/di"
	db_servicetokens "getsturdy.com/api/pkg/servicetokens/db"
)

func Module(c *di.Container) {
	c.Import(db_servicetokens.Module)
	c.Register(New)
}
