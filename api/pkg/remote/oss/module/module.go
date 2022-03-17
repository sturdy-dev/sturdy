package module

import (
	"getsturdy.com/api/pkg/di"
	"getsturdy.com/api/pkg/remote/oss/graphql"
)

func Module(c *di.Container) {
	c.Import(graphql.Module)
}
