package module

import (
	"getsturdy.com/api/pkg/comments/db"
	"getsturdy.com/api/pkg/comments/graphql"
	"getsturdy.com/api/pkg/comments/service"
	"getsturdy.com/api/pkg/di"
)

func Module(c *di.Container) {
	c.Import(db.Module)
	c.Import(service.Module)
	c.Import(graphql.Module)
}
