//go:build !cloud
// +build !cloud

package module

import (
	"getsturdy.com/api/pkg/di"
	"getsturdy.com/api/pkg/installations/db"
	"getsturdy.com/api/pkg/installations/global"
	"getsturdy.com/api/pkg/installations/graphql"
)

func Module(c *di.Container) {
	c.Import(db.Module)
	c.Import(global.Module)
	c.Import(graphql.Module)
}
