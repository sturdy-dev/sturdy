package module

import (
	"mash/pkg/di"
	"mash/pkg/mutagen/db"
)

func Module(c *di.Container) {
	c.Import(db.Module)
}
