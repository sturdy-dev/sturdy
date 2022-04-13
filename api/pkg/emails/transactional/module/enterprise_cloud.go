//go:build enterprise || cloud
// +build enterprise cloud

package module

import (
	"getsturdy.com/api/pkg/di"
	"getsturdy.com/api/pkg/emails/transactional/enterprise"
)

func Module(c *di.Container) {
	c.Import(enterprise.Module)
}
