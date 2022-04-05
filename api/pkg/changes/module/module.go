package module

import (
	"getsturdy.com/api/pkg/changes/db"
	"getsturdy.com/api/pkg/changes/graphql"
	"getsturdy.com/api/pkg/changes/service"
	"getsturdy.com/api/pkg/di"
)

func Module(c *di.Container) {
	c.Import(db.Module)
	c.Import(service.Module)
	c.Import(graphql.Module)
}
