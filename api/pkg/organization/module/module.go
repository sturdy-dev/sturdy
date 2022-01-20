package module

import (
	"getsturdy.com/api/pkg/di"
	"getsturdy.com/api/pkg/organization/db"
	"getsturdy.com/api/pkg/organization/graphql"
	"getsturdy.com/api/pkg/organization/service"
)

func Module(c *di.Container) {
	c.Import(db.Module)
	c.Import(graphql.Module)
	c.Import(service.Module)
}
