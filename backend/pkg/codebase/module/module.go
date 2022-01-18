package module

import (
	"mash/pkg/codebase/db"
	"mash/pkg/codebase/graphql"
	"mash/pkg/codebase/service"
	"mash/pkg/di"
)

func Module(c *di.Container) {
	c.Import(db.Module)
	c.Import(service.Module)
	c.Import(graphql.Module)
}
