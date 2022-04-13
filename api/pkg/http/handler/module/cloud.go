//go:build cloud
// +build cloud

package http

import (
	"getsturdy.com/api/pkg/di"
	"getsturdy.com/api/pkg/http/handler/enterprise/cloud"
)

func Module(c *di.Container) {
	c.Import(cloud.Module)
}
