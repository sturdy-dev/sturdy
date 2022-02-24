//go:build !enterprise && !cloud
// +build !enterprise,!cloud

package module

import (
	"getsturdy.com/api/pkg/di"
	avatars_module "getsturdy.com/api/pkg/users/avatars/module"
	"getsturdy.com/api/pkg/users/db"
	"getsturdy.com/api/pkg/users/graphql"
	oss_selfhosted_service "getsturdy.com/api/pkg/users/oss/selfhosted/service"
	"getsturdy.com/api/pkg/users/service"
)

func Module(c *di.Container) {
	c.Import(db.Module)
	c.Import(graphql.Module)
	c.Import(service.Module)
	c.Register(oss_selfhosted_service.New, new(service.Service))
	c.Import(avatars_module.Module)
}
