package cloud

import (
	"getsturdy.com/api/pkg/di"
	"getsturdy.com/api/pkg/logger/enterprise/cloud/sentry"
)

func Module(c *di.Container) {
	c.Register(sentry.NewClient)
}
