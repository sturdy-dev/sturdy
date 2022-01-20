package module

import (
	"mash/pkg/di"
	"mash/pkg/review/db"
	"mash/pkg/review/graphql"
)

func Module(c *di.Container) {
	c.Import(db.Module)
	c.Import(graphql.Module)
}
