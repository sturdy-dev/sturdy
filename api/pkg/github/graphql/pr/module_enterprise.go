//go:build enterprise || cloud
// +build enterprise cloud

package pr

import (
	"getsturdy.com/api/pkg/di"
	"getsturdy.com/api/pkg/github/graphql/pr/enterprise"
)

func Module(c *di.Container) {
	c.Register(enterprise.NewResolver)
}
