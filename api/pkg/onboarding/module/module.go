package module

import (
	"getsturdy.com/api/pkg/di"
	"getsturdy.com/api/pkg/onboarding/db"
	"getsturdy.com/api/pkg/onboarding/graphql"
)

func Module(c *di.Container) {
	c.Import(db.Module)
	c.Import(graphql.Module)
}
