package module

import (
	"mash/pkg/di"
	"mash/pkg/newsletter/db"
)

func Module(c *di.Container) {
	c.Import(db.Module)
}
