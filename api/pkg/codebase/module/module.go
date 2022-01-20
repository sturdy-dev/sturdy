package module

import (
	"getsturdy.com/api/pkg/codebase/db"
	"getsturdy.com/api/pkg/codebase/graphql"
	"getsturdy.com/api/pkg/codebase/service"
	"getsturdy.com/api/pkg/di"
)

func Module(c *di.Container) {
	c.Import(db.Module)
	c.Import(service.Module)
	c.Import(graphql.Module)
}
