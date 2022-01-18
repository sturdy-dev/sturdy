package module

import (
	"mash/pkg/di"
	"mash/pkg/servicetokens/db"
	"mash/pkg/servicetokens/graphql"
	"mash/pkg/servicetokens/service"
)

func Module(c *di.Container) {
	c.Import(db.Module)
	c.Import(graphql.Module)
	c.Import(service.Module)
}
