//go:build enterprise || cloud
// +build enterprise cloud

package pr

import (
	"mash/pkg/di"
	"mash/pkg/github/graphql/pr/enterprise"
)

func Module(c *di.Container) {
	c.Register(enterprise.NewResolver)
}
