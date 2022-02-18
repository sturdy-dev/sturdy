package selfhosted

import (
	"getsturdy.com/api/pkg/di"
	"getsturdy.com/api/pkg/users/oss/selfhosted/service"
)

func Module(c *di.Container) {
	c.Import(service.Module)
}
