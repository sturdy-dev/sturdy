package module

import (
	"getsturdy.com/api/pkg/di"
	"getsturdy.com/api/pkg/file/graphql"
	"getsturdy.com/api/pkg/file/routes"
	"getsturdy.com/api/pkg/file/service"
)

func Module(c *di.Container) {
	c.Import(graphql.Module)
	c.Import(routes.Module)
	c.Import(service.Module)
}
