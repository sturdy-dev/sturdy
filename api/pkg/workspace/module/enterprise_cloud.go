//go:build enterprise || cloud
// +build enterprise cloud

package module

import (
	"getsturdy.com/api/pkg/di"
	"getsturdy.com/api/pkg/workspace/db"
	"getsturdy.com/api/pkg/workspace/enterprise"
	"getsturdy.com/api/pkg/workspace/graphql"
	"getsturdy.com/api/pkg/workspace/meta"
	"getsturdy.com/api/pkg/workspace/service"
)

func Module(c *di.Container) {
	c.Register(service.New)
	c.Import(graphql.Module)
	c.Import(enterprise.Module)
	c.Import(db.Module)
	c.Import(meta.Module)
}
