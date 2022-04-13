//go:build cloud
// +build cloud

package module

import (
	"getsturdy.com/api/pkg/di"
	cloud "getsturdy.com/api/pkg/users/enterprise/cloud/service"
)

func Module(c *di.Container) {
	c.Import(cloud.Module)
}
