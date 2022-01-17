package module

import (
	"mash/pkg/analytics/configurable"
	"mash/pkg/di"
)

func Module(c *di.Container) {
	c.Register(configurable.NewClient)
}
