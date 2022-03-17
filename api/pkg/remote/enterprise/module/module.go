package module

import (
	"getsturdy.com/api/pkg/di"
	"getsturdy.com/api/pkg/remote/enterprise/db"
	"getsturdy.com/api/pkg/remote/enterprise/graphql"
	"getsturdy.com/api/pkg/remote/enterprise/service"
)

func Module(c *di.Container) {
	c.Import(db.Module)
	c.Import(graphql.Module)
	c.Import(service.Module)
}
