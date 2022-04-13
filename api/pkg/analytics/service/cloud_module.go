//go:build cloud
// +build cloud

package service

import (
	"getsturdy.com/api/pkg/analytics/enterprise/cloud/posthog"
	"getsturdy.com/api/pkg/di"
	"getsturdy.com/api/pkg/logger"
)

func Module(c *di.Container) {
	c.Import(logger.Module)
	c.Import(posthog.Module)
	c.Register(New)
}
