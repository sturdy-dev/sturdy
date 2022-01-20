//go:build !cloud
// +build !cloud

package module

import (
	"mash/pkg/di"
	"mash/pkg/installations/db"
	"mash/pkg/installations/global"
)

func Module(c *di.Container) {
	c.Import(db.Module)
	c.Import(global.Module)
}
