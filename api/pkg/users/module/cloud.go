//go:build cloud
// +build cloud

package module

import (
	"getsturdy.com/api/pkg/di"
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
}
