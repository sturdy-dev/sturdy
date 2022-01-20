//go:build cloud
// +build cloud

package module

import (
	"mash/pkg/analytics/cloud"
	"mash/pkg/di"
)

func Module(c *di.Container) {
	c.Import(cloud.Module)
}
