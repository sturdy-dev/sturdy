package routes

import (
	service_auth "getsturdy.com/api/pkg/auth/service"
	service_change "getsturdy.com/api/pkg/changes/service"
	"getsturdy.com/api/pkg/di"
	service_file "getsturdy.com/api/pkg/file/service"
	"getsturdy.com/api/pkg/logger"
	service_workspaces "getsturdy.com/api/pkg/workspaces/service"
)

func Module(c *di.Container) {
	c.Import(logger.Module)
	c.Import(service_change.Module)
	c.Import(service_workspaces.Module)
	c.Import(service_file.Module)
	c.Import(service_auth.Module)
	c.Register(NewGetFileRoute)
}
