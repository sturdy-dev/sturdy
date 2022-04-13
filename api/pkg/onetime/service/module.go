package service

import (
	"getsturdy.com/api/pkg/di"
	db_onetime "getsturdy.com/api/pkg/onetime/db"
)

func Module(c *di.Container) {
	c.Import(db_onetime.Module)
	c.Register(New)
}
