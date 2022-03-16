//go:build enterprise || cloud
// +build enterprise cloud

package module

import (
	"getsturdy.com/api/pkg/di"
	gh_enterprise "getsturdy.com/api/pkg/github/enterprise"
)

func Module(c *di.Container) {
	c.Import(gh_enterprise.Module)
}
