package testmodule

import (
	"context"

	"github.com/google/uuid"

	module_analytics "getsturdy.com/api/pkg/analytics/module"
	module_configuration "getsturdy.com/api/pkg/configuration/module"
	"getsturdy.com/api/pkg/di"
	module_events "getsturdy.com/api/pkg/events"
	module_events_v2 "getsturdy.com/api/pkg/events/v2"
	"getsturdy.com/api/pkg/installations"
	"getsturdy.com/api/pkg/internal/inmemory"
	module_logger "getsturdy.com/api/pkg/logger/module"
	module_snapshots "getsturdy.com/api/pkg/snapshots/module"
	db_workspaces "getsturdy.com/api/pkg/workspaces/db"
	module_vcs "getsturdy.com/api/vcs/module"
)

func TestModule(c *di.Container) {
	ctx := context.Background()
	c.Register(func() context.Context {
		return ctx
	})

	c.Register(func() *installations.Installation {
		return &installations.Installation{ID: uuid.NewString()}
	})

	c.Import(module_logger.Module)
	c.Import(module_vcs.Module)
	c.Import(module_configuration.TestingModule)
	c.Import(module_analytics.Module)
	c.Import(inmemory.TestModule)
	c.Import(module_snapshots.InMemoryTestingModule)
	c.Import(module_events_v2.Module)
	c.Import(module_events.Module)

	c.Register(func(repo db_workspaces.Repository) (db_workspaces.WorkspaceReader, db_workspaces.WorkspaceWriter) {
		return repo, repo
	})
}
