//go:build cloud
// +build cloud

package module

import (
	"getsturdy.com/api/pkg/di"
	avatars_module "getsturdy.com/api/pkg/users/avatars/module"
	"getsturdy.com/api/pkg/users/db"
	"getsturdy.com/api/pkg/users/enterprise/cloud"
	"getsturdy.com/api/pkg/users/graphql"
	"getsturdy.com/api/pkg/users/service"
)

func Module(c *di.Container) {
	c.Import(db.Module)
	c.Register(service.New)
	c.Register(graphql.NewResolver)
	c.Import(cloud.Module)
	c.Import(avatars_module.Module)
}
