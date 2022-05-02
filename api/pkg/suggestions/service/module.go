package service

import (
	service_analytics "getsturdy.com/api/pkg/analytics/service"
	service_codebases "getsturdy.com/api/pkg/codebases/service"
	"getsturdy.com/api/pkg/di"
	"getsturdy.com/api/pkg/events"
	"getsturdy.com/api/pkg/logger"
	sender_notification "getsturdy.com/api/pkg/notification/sender"
	service_snapshots "getsturdy.com/api/pkg/snapshots/service"
	db_suggestions "getsturdy.com/api/pkg/suggestions/db"
	service_worspace "getsturdy.com/api/pkg/workspaces/service"
	"getsturdy.com/api/vcs/executor"
)

func Module(c *di.Container) {
	c.Import(logger.Module)
	c.Import(db_suggestions.Module)
	c.Import(service_worspace.Module)
	c.Import(executor.Module)
	c.Import(service_snapshots.Module)
	c.Import(service_analytics.Module)
	c.Import(sender_notification.Module)
	c.Import(service_codebases.Module)
	c.Import(events.Module)
	c.Register(New)
}
