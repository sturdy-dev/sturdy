package webhooks

import (
	sender_workspace_activity "getsturdy.com/api/pkg/activity/sender"
	service_activity "getsturdy.com/api/pkg/activity/service"
	service_analytics "getsturdy.com/api/pkg/analytics/service"
	service_change "getsturdy.com/api/pkg/changes/service"
	workers_ci "getsturdy.com/api/pkg/ci/workers"
	db_codebases "getsturdy.com/api/pkg/codebases/db"
	service_codebases "getsturdy.com/api/pkg/codebases/service"
	service_comments "getsturdy.com/api/pkg/comments/service"
	configuration "getsturdy.com/api/pkg/configuration/module"
	"getsturdy.com/api/pkg/di"
	"getsturdy.com/api/pkg/events"
	eventsv2 "getsturdy.com/api/pkg/events/v2"
	github_client "getsturdy.com/api/pkg/github/enterprise/client"
	db_github "getsturdy.com/api/pkg/github/enterprise/db"
	service_github "getsturdy.com/api/pkg/github/enterprise/service"
	service_github_importing "getsturdy.com/api/pkg/github/enterprise/service/importing"
	"getsturdy.com/api/pkg/logger"
	db_review "getsturdy.com/api/pkg/review/db"
	service_statuses "getsturdy.com/api/pkg/statuses/service"
	service_sync "getsturdy.com/api/pkg/sync/service"
	service_users "getsturdy.com/api/pkg/users/service/module"
	db_view "getsturdy.com/api/pkg/views/db"
	db_workspaces "getsturdy.com/api/pkg/workspaces/db"
	meta_workspaces "getsturdy.com/api/pkg/workspaces/meta"
	"getsturdy.com/api/vcs/executor"
)

func Module(c *di.Container) {
	c.Import(logger.Module)
	c.Import(configuration.Module)
	c.Import(github_client.Module)
	c.Import(db_github.Module)
	c.Import(db_workspaces.Module)
	c.Import(db_review.Module)
	c.Import(db_codebases.Module)
	c.Import(db_view.Module)
	c.Import(service_analytics.Module)
	c.Import(service_sync.Module)
	c.Import(service_codebases.Module)
	c.Import(service_comments.Module)
	c.Import(service_activity.Module)
	c.Import(service_change.Module)
	c.Import(service_statuses.Module)
	c.Import(service_github.Module)
	c.Import(service_github_importing.Module)
	c.Import(service_users.Module)
	c.Import(workers_ci.Module)
	c.Import(eventsv2.Module)
	c.Import(events.Module)
	c.Import(sender_workspace_activity.Module)
	c.Import(meta_workspaces.Module)
	c.Import(executor.Module)
	c.Register(New)
	c.Register(NewWebhooksQueue)
}
