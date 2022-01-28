package selfhosted

import (
	"getsturdy.com/api/pkg/di"
	"getsturdy.com/api/pkg/features/enterprise/selfhosted/graphql"
)

func Module(c *di.Container) {
	c.Import(graphql.Module)
}
