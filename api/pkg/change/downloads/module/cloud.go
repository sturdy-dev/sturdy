//go:build cloud
// +build cloud

package module

import (
	"getsturdy.com/api/pkg/change/downloads/enterprise/cloud"
	"getsturdy.com/api/pkg/di"
)

func Module(c *di.Container) {
	c.Import(cloud.Module)
}
