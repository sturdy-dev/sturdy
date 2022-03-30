package module

import (
	"getsturdy.com/api/pkg/di"
	"getsturdy.com/api/pkg/remote/oss/graphql"
	"getsturdy.com/api/pkg/remote/oss/service"
)

func Module(c *di.Container) {
	c.Import(graphql.Module)
	c.Import(service.Module)
}
