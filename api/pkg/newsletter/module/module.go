package module

import (
	"getsturdy.com/api/pkg/di"
	"getsturdy.com/api/pkg/newsletter/db"
)

func Module(c *di.Container) {
	c.Import(db.Module)
}
