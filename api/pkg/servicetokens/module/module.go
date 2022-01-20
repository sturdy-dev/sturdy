package module

import (
	"getsturdy.com/api/pkg/di"
	"getsturdy.com/api/pkg/servicetokens/db"
	"getsturdy.com/api/pkg/servicetokens/graphql"
	"getsturdy.com/api/pkg/servicetokens/service"
)

func Module(c *di.Container) {
	c.Import(db.Module)
	c.Import(graphql.Module)
	c.Import(service.Module)
}
