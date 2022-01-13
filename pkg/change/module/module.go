package module

import (
	"mash/pkg/change/db"
	"mash/pkg/change/graphql"
	"mash/pkg/change/service"
	"mash/pkg/di"
)

func Module(c *di.Container) {
	c.Import(db.Module)
	c.Import(service.Module)
	c.Import(graphql.Module)
}
