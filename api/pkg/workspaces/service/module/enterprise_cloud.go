//go:build enterprise || cloud
// +build enterprise cloud

package module

import (
	"getsturdy.com/api/pkg/di"
	enterprise "getsturdy.com/api/pkg/workspaces/enterprise/service"
)

func Module(c *di.Container) {
	c.Import(enterprise.Module)
}
