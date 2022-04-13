package service

import (
	"getsturdy.com/api/pkg/codebases/acl/provider"
	"getsturdy.com/api/pkg/di"
	db_snapshots "getsturdy.com/api/pkg/snapshots/db"
	"getsturdy.com/api/vcs/executor"
)

func Module(c *di.Container) {
	c.Import(executor.Module)
	c.Import(provider.Module)
	c.Import(db_snapshots.Module)
	c.Register(New)
}
