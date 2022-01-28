//go:build !enterprise
// +build !enterprise

package module

import (
	"getsturdy.com/api/pkg/di"
	"getsturdy.com/api/pkg/user/db"
	"getsturdy.com/api/pkg/user/graphql"
	"getsturdy.com/api/pkg/user/service"
)

func Module(c *di.Container) {
	c.Import(db.Module)
	c.Import(graphql.Module)
	c.Register(service.New, new(service.Service))
}
