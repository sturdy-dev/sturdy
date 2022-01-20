package module

import (
	"getsturdy.com/api/pkg/di"
	"getsturdy.com/api/pkg/features/graphql"
)

func Module(c *di.Container) {
	c.Import(graphql.Module)
}
