package module

import (
	"mash/pkg/di"
	"mash/pkg/presence/db"
	"mash/pkg/presence/graphql"
	"mash/pkg/presence/service"
)

func Module(c *di.Container) {
	c.Import(db.Module)
	c.Import(graphql.Module)
	c.Import(service.Module)
}
