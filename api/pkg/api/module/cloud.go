//go:build cloud
// +build cloud

package api

import (
	"getsturdy.com/api/pkg/api/enterprise/cloud"
	"getsturdy.com/api/pkg/di"
)

func Module(c *di.Container) {
	c.Import(cloud.Module)
}
