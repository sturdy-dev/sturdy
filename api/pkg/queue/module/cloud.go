//go:build cloud
// +build cloud

package queue

import (
	"getsturdy.com/api/pkg/di"
	"getsturdy.com/api/pkg/queue/enterprise/cloud"
)

func Module(c *di.Container) {
	c.Import(cloud.Module)
}
