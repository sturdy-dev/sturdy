//go:build cloud
// +build cloud

package aws

import (
	"getsturdy.com/api/pkg/aws/enterprise/cloud"
	"getsturdy.com/api/pkg/di"
)

func Module(c *di.Container) {
	c.Import(cloud.Module)
}
