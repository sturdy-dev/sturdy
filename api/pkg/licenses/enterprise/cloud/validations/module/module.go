package module

import (
	"getsturdy.com/api/pkg/di"
	"getsturdy.com/api/pkg/licenses/enterprise/cloud/validations/db"
	"getsturdy.com/api/pkg/licenses/enterprise/cloud/validations/service"
)

func Module(c *di.Container) {
	c.Import(db.Module)
	c.Import(service.Module)
}
