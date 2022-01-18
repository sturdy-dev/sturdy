package cloud

import (
	"mash/pkg/analytics/cloud/posthog"
	"mash/pkg/di"
)

func Module(c *di.Container) {
	c.Import(posthog.Module)
}
