package service

import (
	sender_activity "getsturdy.com/api/pkg/activity/sender"
	service_activity "getsturdy.com/api/pkg/activity/service"
	service_analytics "getsturdy.com/api/pkg/analytics/service"
	service_change "getsturdy.com/api/pkg/changes/service"
	workers_ci "getsturdy.com/api/pkg/ci/workers"
	db_codebases "getsturdy.com/api/pkg/codebases/db"
	service_comments "getsturdy.com/api/pkg/comments/service"
	configuration "getsturdy.com/api/pkg/configuration/module"
	"getsturdy.com/api/pkg/di"
	emails "getsturdy.com/api/pkg/emails/module"
	"getsturdy.com/api/pkg/events"
	eventsv2 "getsturdy.com/api/pkg/events/v2"
	"getsturdy.com/api/pkg/github/enterprise/client"
	db_github "getsturdy.com/api/pkg/github/enterprise/db"
	"getsturdy.com/api/pkg/logger"
	sender_notification "getsturdy.com/api/pkg/notification/sender"
	service_remote "getsturdy.com/api/pkg/remote/enterprise/service"
	db_review "getsturdy.com/api/pkg/review/db"
	"getsturdy.com/api/pkg/snapshots/snapshotter"
	service_sync "getsturdy.com/api/pkg/sync/service"
	service_user "getsturdy.com/api/pkg/users/service/module"
	db_workspaces "getsturdy.com/api/pkg/workspaces/db"
	meta_workspaces "getsturdy.com/api/pkg/workspaces/meta"
	service_workspaces "getsturdy.com/api/pkg/workspaces/service"
	"getsturdy.com/api/vcs/executor"
)

func Module(c *di.Container) {
	c.Import(logger.Module)
	c.Import(db_github.Module)
	c.Import(configuration.Module)
	c.Import(client.Module)
	c.Import(db_workspaces.Module)
	c.Import(meta_workspaces.Module)
	c.Import(db_codebases.Module)
	c.Import(db_review.Module)
	c.Import(executor.Module)
	c.Import(snapshotter.Module)
	c.Import(service_analytics.Module)
	c.Import(emails.Module)
	c.Import(sender_notification.Module)
	c.Import(sender_activity.Module)
	c.Import(events.Module)
	c.Import(eventsv2.Module)
	c.Import(service_user.Module)
	c.Import(service_workspaces.Module)
	c.Import(service_sync.Module)
	c.Import(service_comments.Module)
	c.Import(service_change.Module)
	c.Import(service_remote.Module)
	c.Import(workers_ci.Module)
	c.Import(service_activity.Module)
	c.Register(NewClonerQueue)
	c.Register(NewImporterQueue)
	c.Register(New)
}
