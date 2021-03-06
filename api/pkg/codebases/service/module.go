package service

import (
	sender_notifications "getsturdy.com/api/pkg/activity/sender"
	service_analytics "getsturdy.com/api/pkg/analytics/service"
	service_changes "getsturdy.com/api/pkg/changes/service"
	db_codebases "getsturdy.com/api/pkg/codebases/db"
	"getsturdy.com/api/pkg/di"
	"getsturdy.com/api/pkg/events"
	"getsturdy.com/api/pkg/logger"
	service_users "getsturdy.com/api/pkg/users/service/module"
	service_workspaces "getsturdy.com/api/pkg/workspaces/service"
	"getsturdy.com/api/vcs/executor"
)

func Module(c *di.Container) {
	c.Import(db_codebases.Module)
	c.Import(service_workspaces.Module)
	c.Import(service_users.Module)
	c.Import(logger.Module)
	c.Import(executor.Module)
	c.Import(events.Module)
	c.Import(service_analytics.Module)
	c.Import(service_changes.Module)
	c.Import(sender_notifications.Module)
	c.Register(New)
}
