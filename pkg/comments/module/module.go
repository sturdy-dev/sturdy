package module

import (
	"mash/pkg/comments/db"
	"mash/pkg/comments/graphql"
	"mash/pkg/comments/service"
	"mash/pkg/di"
)

func Module(c *di.Container) {
	c.Import(db.Module)
	c.Import(service.Module)
	c.Import(graphql.Module)
}
