//go:build !cloud && !enterprise
// +build !cloud,!enterprise

package module

import (
	"getsturdy.com/api/pkg/di"
	"getsturdy.com/api/pkg/installations/db"
	"getsturdy.com/api/pkg/installations/global"
	"getsturdy.com/api/pkg/installations/oss/graphql"
	"getsturdy.com/api/pkg/installations/service"
)

func Module(c *di.Container) {
	c.Import(db.Module)
	c.Import(global.Module)
	c.Import(graphql.Module)
	c.Import(service.Module)
}
