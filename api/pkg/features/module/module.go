package module

import (
	"mash/pkg/di"
	"mash/pkg/features/graphql"
)

func Module(c *di.Container) {
	c.Import(graphql.Module)
}
