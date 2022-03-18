package module

import (
	"getsturdy.com/api/pkg/codebases/db"
	"getsturdy.com/api/pkg/codebases/graphql"
	"getsturdy.com/api/pkg/codebases/service"
	"getsturdy.com/api/pkg/di"
)

func Module(c *di.Container) {
	c.Import(db.Module)
	c.Import(service.Module)
	c.Import(graphql.Module)
}
