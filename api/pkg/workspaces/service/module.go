package service

import (
	"getsturdy.com/api/pkg/activity/sender"
	service_activity "getsturdy.com/api/pkg/activity/service"
	service_analytics "getsturdy.com/api/pkg/analytics/service"
	service_changes "getsturdy.com/api/pkg/changes/service"
	"getsturdy.com/api/pkg/di"
	"getsturdy.com/api/pkg/events"
	eventsv2 "getsturdy.com/api/pkg/events/v2"
	"getsturdy.com/api/pkg/logger"
	service_snapshots "getsturdy.com/api/pkg/snapshots/service"
	service_users "getsturdy.com/api/pkg/users/service/module"
	service_view "getsturdy.com/api/pkg/view/service"
	db_workspaces "getsturdy.com/api/pkg/workspaces/db"
	"getsturdy.com/api/vcs/executor"
)

func Module(c *di.Container) {
	c.Import(logger.Module)
	c.Import(service_analytics.Module)
	c.Import(db_workspaces.Module)
	c.Import(service_view.Module)
	c.Import(service_activity.Module)
	c.Import(service_users.Module)
	c.Import(service_changes.Module)
	c.Import(sender.Module)
	c.Import(executor.Module)
	c.Import(events.Module)
	c.Import(eventsv2.Module)
	c.Import(service_snapshots.Module)
	c.Register(New)
}
