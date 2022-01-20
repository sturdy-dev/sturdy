package module

import (
	"getsturdy.com/api/pkg/di"
	"getsturdy.com/api/pkg/presence/db"
	"getsturdy.com/api/pkg/presence/graphql"
	"getsturdy.com/api/pkg/presence/service"
)

func Module(c *di.Container) {
	c.Import(db.Module)
	c.Import(graphql.Module)
	c.Import(service.Module)
}
