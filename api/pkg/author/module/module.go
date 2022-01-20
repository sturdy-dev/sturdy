package module

import (
	"mash/pkg/author/graphql"
	"mash/pkg/di"
)

func Module(c *di.Container) {
	c.Import(graphql.Module)
}
