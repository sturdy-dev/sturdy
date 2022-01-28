package module

import (
	"getsturdy.com/api/pkg/di"
	"getsturdy.com/api/pkg/user/db"
	"getsturdy.com/api/pkg/user/graphql"
	"getsturdy.com/api/pkg/user/service"
)

func Module(c *di.Container) {
	c.Import(service.Module)
	c.Import(db.Module)
	c.Import(graphql.Module)
}
