package service

import (
	service_changes "getsturdy.com/api/pkg/changes/service"
	provider_acl "getsturdy.com/api/pkg/codebases/acl/provider"
	service_codebases "getsturdy.com/api/pkg/codebases/service"
	"getsturdy.com/api/pkg/di"
	service_organizations "getsturdy.com/api/pkg/organization/service"
	service_users "getsturdy.com/api/pkg/users/service"
	service_workspaces "getsturdy.com/api/pkg/workspaces/service"
)

func Module(c *di.Container) {
	c.Import(service_codebases.Module)
	c.Import(service_changes.Module)
	c.Import(service_users.Module)
	c.Import(service_workspaces.Module)
	c.Import(service_organizations.Module)
	c.Import(provider_acl.Module)
	c.Register(New)
}
