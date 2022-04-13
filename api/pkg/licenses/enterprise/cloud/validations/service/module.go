package service

import (
	"getsturdy.com/api/pkg/di"
	"getsturdy.com/api/pkg/licenses/enterprise/cloud/validations/db"
)

func Module(c *di.Container) {
	c.Import(db.Module)
	c.Register(New)
}
