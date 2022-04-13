package snapshotter

import (
	service_analytics "getsturdy.com/api/pkg/analytics/service"
	"getsturdy.com/api/pkg/di"
	"getsturdy.com/api/pkg/events"
	eventsv2 "getsturdy.com/api/pkg/events/v2"
	"getsturdy.com/api/pkg/logger"
	db_snapshots "getsturdy.com/api/pkg/snapshots/db"
	db_suggestions "getsturdy.com/api/pkg/suggestions/db"
	db_view "getsturdy.com/api/pkg/view/db"
	db_workspaces "getsturdy.com/api/pkg/workspaces/db"
	meta_workspaces "getsturdy.com/api/pkg/workspaces/meta"
	"getsturdy.com/api/vcs/executor"
)

func Module(c *di.Container) {
	c.Import(db_snapshots.Module)
	c.Import(db_workspaces.Module)
	c.Import(meta_workspaces.Module)
	c.Import(db_view.Module)
	c.Import(events.Module)
	c.Import(eventsv2.Module)
	c.Import(logger.Module)
	c.Import(db_suggestions.Module)
	c.Import(executor.Module)
	c.Import(service_analytics.Module)
	c.Register(NewGitSnapshotter)
}
