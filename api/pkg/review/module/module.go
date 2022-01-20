package module

import (
	"getsturdy.com/api/pkg/di"
	"getsturdy.com/api/pkg/review/db"
	"getsturdy.com/api/pkg/review/graphql"
)

func Module(c *di.Container) {
	c.Import(db.Module)
	c.Import(graphql.Module)
}
