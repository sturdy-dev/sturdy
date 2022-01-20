package module

import (
	"mash/pkg/di"
	"mash/pkg/onboarding/db"
	"mash/pkg/onboarding/graphql"
)

func Module(c *di.Container) {
	c.Import(db.Module)
	c.Import(graphql.Module)
}
