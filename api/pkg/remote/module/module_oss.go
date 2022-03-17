//go:build !enterprise && !cloud

package module

import (
	"getsturdy.com/api/pkg/di"
	oss_module "getsturdy.com/api/pkg/remote/oss/module"
)

func Module(c *di.Container) {
	c.Import(oss_module.Module)
}
