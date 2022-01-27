package cloud

import (
	"getsturdy.com/api/pkg/di"
	"getsturdy.com/api/pkg/licenses/enterprise/cloud/db"
	"getsturdy.com/api/pkg/licenses/enterprise/cloud/graphql"
	"getsturdy.com/api/pkg/licenses/enterprise/cloud/service"
	module_validations "getsturdy.com/api/pkg/licenses/enterprise/cloud/validations/module"
)

func Module(c *di.Container) {
	c.Import(module_validations.Module)
	c.Import(db.Module)
	c.Import(service.Module)
	c.Import(graphql.Module)
}
