package graphql

import (
	service_auth "getsturdy.com/api/pkg/auth/service"
	service_codebase "getsturdy.com/api/pkg/codebases/service"
	graphql_crypto "getsturdy.com/api/pkg/crypto/graphql"
	"getsturdy.com/api/pkg/di"
	"getsturdy.com/api/pkg/remote/enterprise/service"
	service_user "getsturdy.com/api/pkg/users/service/module"
	service_workspace "getsturdy.com/api/pkg/workspaces/service"
)

func Module(c *di.Container) {
	c.Import(service.Module)
	c.Import(service_workspace.Module)
	c.Import(service_auth.Module)
	c.Import(service_codebase.Module)
	c.Import(service_user.Module)
	c.Import(graphql_crypto.Module)
	c.Register(New)
}
