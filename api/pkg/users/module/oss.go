//go:build !enterprise && !cloud
// +build !enterprise,!cloud

package module

import (
	"getsturdy.com/api/pkg/di"
	"getsturdy.com/api/pkg/users/db"
	"getsturdy.com/api/pkg/users/graphql"
	"getsturdy.com/api/pkg/users/service"
)

func Module(c *di.Container) {
	c.Import(db.Module)
	c.Import(graphql.Module)
	c.Register(service.New, new(service.Service))
}
