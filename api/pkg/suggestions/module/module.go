package module

import (
	"getsturdy.com/api/pkg/di"
	"getsturdy.com/api/pkg/suggestions/db"
	"getsturdy.com/api/pkg/suggestions/graphql"
	"getsturdy.com/api/pkg/suggestions/service"
)

func Module(c *di.Container) {
	c.Import(db.Module)
	c.Import(graphql.Module)
	c.Import(service.Module)
}
