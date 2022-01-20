package module

import (
	"mash/pkg/di"
	"mash/pkg/notification/db"
	"mash/pkg/notification/graphql"
	"mash/pkg/notification/sender"
	"mash/pkg/notification/service"
)

func Module(c *di.Container) {
	c.Import(db.Module)
	c.Import(graphql.Module)
	c.Import(sender.Module)
	c.Import(service.Module)
}
