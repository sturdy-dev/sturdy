package module

import (
	"mash/pkg/di"
	"mash/pkg/organization/db"
	"mash/pkg/organization/graphql"
	"mash/pkg/organization/service"
)

func Module(c *di.Container) {
	c.Import(db.Module)
	c.Import(graphql.Module)
	c.Import(service.Module)
}
