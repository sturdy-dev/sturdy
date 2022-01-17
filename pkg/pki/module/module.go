package module

import (
	"mash/pkg/di"
	"mash/pkg/pki/db"
	"mash/pkg/pki/graphql"
)

func Module(c *di.Container) {
	c.Import(db.Module)
	c.Import(graphql.Module)
}
