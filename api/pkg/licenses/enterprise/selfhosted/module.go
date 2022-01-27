package selfhosted

import (
	"getsturdy.com/api/pkg/di"
	"getsturdy.com/api/pkg/licenses/enterprise/selfhosted/graphql"
	"getsturdy.com/api/pkg/licenses/enterprise/selfhosted/validator"
)

func Module(c *di.Container) {
	c.Import(graphql.Module)
	c.Import(validator.Module)
}
