package module

import (
	"getsturdy.com/api/pkg/di"
	"getsturdy.com/api/pkg/workspaces/watchers/db"
	"getsturdy.com/api/pkg/workspaces/watchers/graphql"
	"getsturdy.com/api/pkg/workspaces/watchers/service"
)

func Module(c *di.Container) {
	c.Import(db.Module)
	c.Import(service.Module)
	c.Import(graphql.Module)
}
