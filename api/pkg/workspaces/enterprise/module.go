package enterprise

import (
	"getsturdy.com/api/pkg/di"
	"getsturdy.com/api/pkg/workspaces/enterprise/service"
)

func Module(c *di.Container) {
	c.Import(service.Module)
}
