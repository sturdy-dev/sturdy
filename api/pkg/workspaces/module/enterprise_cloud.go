//go:build enterprise || cloud
// +build enterprise cloud

package module

import (
	"getsturdy.com/api/pkg/di"
	"getsturdy.com/api/pkg/workspaces/db"
	"getsturdy.com/api/pkg/workspaces/enterprise"
	"getsturdy.com/api/pkg/workspaces/graphql"
	"getsturdy.com/api/pkg/workspaces/meta"
	"getsturdy.com/api/pkg/workspaces/service"
)

func Module(c *di.Container) {
	c.Register(service.New)
	c.Import(graphql.Module)
	c.Import(enterprise.Module)
	c.Import(db.Module)
	c.Import(meta.Module)
}
