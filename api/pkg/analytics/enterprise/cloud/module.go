package cloud

import (
	"getsturdy.com/api/pkg/analytics/enterprise/cloud/posthog"
	"getsturdy.com/api/pkg/di"
)

func Module(c *di.Container) {
	c.Import(posthog.Module)
}
