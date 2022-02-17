//go:build cloud
// +build cloud

package module

import (
	"getsturdy.com/api/pkg/analytics/enterprise/cloud"
	"getsturdy.com/api/pkg/analytics/service"
	"getsturdy.com/api/pkg/di"
)

func Module(c *di.Container) {
	c.Import(cloud.Module)
	c.Import(service.Module)
}
