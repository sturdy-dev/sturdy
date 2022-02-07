//go:build cloud
// +build cloud

package module

import (
	"getsturdy.com/api/pkg/di"
	"getsturdy.com/api/pkg/emails/enterprise/cloud"
)

func Module(c *di.Container) {
	c.Register(cloud.New)
}
