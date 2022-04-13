package service

import (
	db_changes "getsturdy.com/api/pkg/changes/db"
	db_codebases "getsturdy.com/api/pkg/codebases/db"
	"getsturdy.com/api/pkg/di"
	"getsturdy.com/api/pkg/logger"
	"getsturdy.com/api/pkg/snapshots/snapshotter"
	"getsturdy.com/api/vcs/executor"
)

func Module(c *di.Container) {
	c.Import(db_changes.Module)
	c.Import(db_codebases.Module)
	c.Import(executor.Module)
	c.Import(snapshotter.Module)
	c.Import(logger.Module)
	c.Register(New)
}
