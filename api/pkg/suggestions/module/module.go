package module

import (
	"mash/pkg/di"
	"mash/pkg/suggestions/db"
	"mash/pkg/suggestions/graphql"
	"mash/pkg/suggestions/service"
)

func Module(c *di.Container) {
	c.Import(db.Module)
	c.Import(graphql.Module)
	c.Import(service.Module)
}
