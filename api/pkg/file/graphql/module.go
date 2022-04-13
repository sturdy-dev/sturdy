package graphql

import (
	service_auth "getsturdy.com/api/pkg/auth/service"
	service_change "getsturdy.com/api/pkg/changes/service"
	"getsturdy.com/api/pkg/di"
	service_file "getsturdy.com/api/pkg/file/service"
	"getsturdy.com/api/vcs/executor"
)

func Module(c *di.Container) {
	c.Import(executor.Module)
	c.Import(service_auth.Module)
	c.Import(service_change.Module)
	c.Import(service_file.Module)
	c.Register(NewFileRootResolver)
}
