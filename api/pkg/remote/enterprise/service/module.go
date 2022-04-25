package service

import (
	analytics_service "getsturdy.com/api/pkg/analytics/service"
	service_change "getsturdy.com/api/pkg/changes/service"
	db_crypto "getsturdy.com/api/pkg/crypto/db"
	"getsturdy.com/api/pkg/di"
	"getsturdy.com/api/pkg/logger"
	db_remote "getsturdy.com/api/pkg/remote/enterprise/db"
	remote_service "getsturdy.com/api/pkg/remote/service"
	service_snapshots "getsturdy.com/api/pkg/snapshots/service"
	db_workspaces "getsturdy.com/api/pkg/workspaces/db"
	meta_workspaces "getsturdy.com/api/pkg/workspaces/meta"
	"getsturdy.com/api/vcs/executor"
)

func Module(c *di.Container) {
	c.Import(db_remote.Module)
	c.Import(executor.Module)
	c.Import(logger.Module)
	c.Import(db_workspaces.Module)
	c.Import(meta_workspaces.Module)
	c.Import(service_snapshots.Module)
	c.Import(service_change.Module)
	c.Import(analytics_service.Module)
	c.Import(db_crypto.Module)
	c.Register(New)
	c.Register(func(e *EnterpriseService) remote_service.Service {
		return e
	})
}
