//go:build enterprise || cloud
// +build enterprise cloud

package pr

import (
	"getsturdy.com/api/pkg/di"
	enterprise "getsturdy.com/api/pkg/github/enterprise/graphql/pr"
)

func Module(c *di.Container) {
	c.Import(enterprise.Module)
}
