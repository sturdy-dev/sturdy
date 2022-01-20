//go:build enterprise || cloud
// +build enterprise cloud

package module

import (
	"mash/pkg/di"
	"mash/pkg/integrations/buildkite/enterprise"
)

func Module(c *di.Container) {
	c.Import(enterprise.Module)
}
