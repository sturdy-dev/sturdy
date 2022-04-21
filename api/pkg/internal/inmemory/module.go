package inmemory

import (
	"getsturdy.com/api/pkg/di"
	db_suggestions "getsturdy.com/api/pkg/suggestions/db"
	db_workspaces "getsturdy.com/api/pkg/workspaces/db"
)

func TestModule(c *di.Container) {
	c.Register(NewInMemoryAclRepo)
	c.Register(NewInMemoryGitHubInstallationRepository)
	c.Register(NewInMemoryGitHubRepositoryRepo)
	c.Register(NewInMemoryGitHubUserRepo)
	c.Register(NewInMemorySnapshotRepo)
	c.Register(NewInMemoryViewRepo)
	c.Register(db_workspaces.NewMemory)
	c.Register(db_suggestions.NewMemory)
}
