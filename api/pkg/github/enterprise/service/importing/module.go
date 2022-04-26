package importing

import (
	sender_activity "getsturdy.com/api/pkg/activity/sender"
	workers_ci "getsturdy.com/api/pkg/ci/workers"
	service_comments "getsturdy.com/api/pkg/comments/service"
	"getsturdy.com/api/pkg/di"
	db_github "getsturdy.com/api/pkg/github/enterprise/db"
	service_github "getsturdy.com/api/pkg/github/enterprise/service"
	"getsturdy.com/api/pkg/logger"
	db_workspaces "getsturdy.com/api/pkg/workspaces/db"
	service_workspaces "getsturdy.com/api/pkg/workspaces/service"
	"getsturdy.com/api/vcs/executor"
)

func Module(c *di.Container) {
	c.Import(logger.Module)
	c.Import(db_github.Module)
	c.Import(db_workspaces.Module)
	c.Import(executor.Module)
	c.Import(sender_activity.Module)
	c.Import(service_workspaces.Module)
	c.Import(service_comments.Module)
	c.Import(service_github.Module)
	c.Import(workers_ci.Module)
	c.Register(New)
}
