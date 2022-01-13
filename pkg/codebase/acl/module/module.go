package module

import (
	"mash/pkg/codebase/acl/db"
	"mash/pkg/codebase/acl/graphql"
	"mash/pkg/codebase/acl/provider"
	"mash/pkg/di"
)

func Module(c *di.Container) {
	c.Import(provider.Module)
	c.Import(graphql.Module)
	c.Import(db.Module)
}
