//go:build enterprise || cloud
// +build enterprise cloud

package graphql

import (
	"getsturdy.com/api/pkg/di"
	"getsturdy.com/api/pkg/statuses/db"
	"getsturdy.com/api/pkg/statuses/enterprise"
	"getsturdy.com/api/pkg/statuses/graphql"
	"getsturdy.com/api/pkg/statuses/service"
)

func Module(c *di.Container) {
	c.Register(graphql.New)
	c.Import(db.Module)
	c.Import(service.Module)
	c.Import(enterprise.Module)
}
