//go:build !cloud
// +build !cloud

package module

import (
	"getsturdy.com/api/pkg/di"
	"getsturdy.com/api/pkg/installations/db"
	"getsturdy.com/api/pkg/installations/global"
)

func Module(c *di.Container) {
	c.Import(db.Module)
	c.Import(global.Module)
}
