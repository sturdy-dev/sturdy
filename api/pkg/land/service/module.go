package service

import (
	"getsturdy.com/api/pkg/activity/sender"
	service_activity "getsturdy.com/api/pkg/activity/service"
	service_analytics "getsturdy.com/api/pkg/analytics/service"
	service_changes "getsturdy.com/api/pkg/changes/service"
	workers_ci "getsturdy.com/api/pkg/ci/workers"
	service_comments "getsturdy.com/api/pkg/comments/service"
	"getsturdy.com/api/pkg/di"
	"getsturdy.com/api/pkg/events"
	eventsv2 "getsturdy.com/api/pkg/events/v2"
	"getsturdy.com/api/pkg/logger"
	"getsturdy.com/api/pkg/snapshots/snapshotter"
	worker_snapshots "getsturdy.com/api/pkg/snapshots/worker"
	service_users "getsturdy.com/api/pkg/users/service/module"
	service_view "getsturdy.com/api/pkg/view/service"
	db_workspaces "getsturdy.com/api/pkg/workspaces/db"
	service_workspaces "getsturdy.com/api/pkg/workspaces/service"
	"getsturdy.com/api/vcs/executor"
)

func Module(c *di.Container) {
	c.Import(logger.Module)
	c.Import(db_workspaces.Module)
	c.Import(service_users.Module)
	c.Import(service_workspaces.Module)
	c.Import(service_changes.Module)
	c.Import(service_analytics.Module)
	c.Import(service_view.Module)
	c.Import(service_comments.Module)
	c.Import(service_activity.Module)
	c.Import(snapshotter.Module)
	c.Import(worker_snapshots.Module)
	c.Import(events.Module)
	c.Import(eventsv2.Module)
	c.Import(executor.Module)
	c.Import(workers_ci.Module)
	c.Import(sender.Module)
	c.Register(New)
}
