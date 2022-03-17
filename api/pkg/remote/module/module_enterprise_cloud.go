//go:build enterprise || cloud

package module

import (
	"getsturdy.com/api/pkg/di"
	enterprise_module "getsturdy.com/api/pkg/remote/enterprise/module"
)

func Module(c *di.Container) {
	c.Import(enterprise_module.Module)
}
