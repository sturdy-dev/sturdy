package graphql

import (
	service_auth "getsturdy.com/api/pkg/auth/service"
	service_codebase "getsturdy.com/api/pkg/codebases/service"
	"getsturdy.com/api/pkg/di"
	service_servicetokens "getsturdy.com/api/pkg/servicetokens/service"
)

func Module(c *di.Container) {
	c.Import(service_auth.Module)
	c.Import(service_codebase.Module)
	c.Import(service_servicetokens.Module)
	c.Register(New)
}
